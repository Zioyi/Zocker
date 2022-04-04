package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Zioyi/zocker/cgroups"
	"github.com/Zioyi/zocker/cgroups/subsystems"
	"github.com/Zioyi/zocker/container"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string, containerName string) {
	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName)
	if parent == nil {
		log.Errorf("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerID)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	// 启动时限制
	cgroupManager := cgroups.NewCgroupManger("zocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	log.Infof("[zocker] container name is %s\n", containerName)
	if tty {
		_ = parent.Wait()
		mntURL := "/root/mnt/"
		rootURL := "/root/"
		container.DeleteWorkSpace(rootURL, mntURL, volume)
		deleteContainerInfo(containerName)
	}
	os.Exit(0)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName string, id string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, " ")
	containerInfo := &container.ContainerInfo{
		Id:         id,
		Pid:        strconv.Itoa(containerPID),
		Command:    command,
		Status:     container.Running,
		Name:       containerName,
		CreateTime: createTime,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record contaier info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Close %s error %v", fileName, err)
		}
	}(file)
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}

	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containName string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
