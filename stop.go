package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"syscall"

	"github.com/Zioyi/zocker/container"

	log "github.com/sirupsen/logrus"
)

func StopContainer(containerName string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Errorf("[StopContainer] Get container pid by name %s failed, err = %v", containerName, err)
		return
	}

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("[StopContainer] Convert pid from string to int failed, err = %v", err)
		return
	}

	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("[StopContaienr] Kill container %s failed, err = %v", containerName, err)
		return
	}

	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("[StopContainer] Get containerInfo by name %s fialed, err = %v", containerName, err)
		return
	}

	containerInfo.Status = container.Stop
	containerInfo.Pid = ""

	err = UpdateContainerInfo(containerInfo)
	if err != nil {
		log.Errorf("[StopContainer] Update contaienrInfo %s failed, err = %v", containerName, err)
	}

	return
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var containerInfo container.ContainerInfo

	if err = json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return nil, err
	}

	return &containerInfo, nil
}

func UpdateContainerInfo(containerInfo *container.ContainerInfo) error {
	contentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return err
	}

	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerInfo.Name)
	configFilePath := dirURL + container.ConfigName
	err = ioutil.WriteFile(configFilePath, contentBytes, 0622)

	return err

}
