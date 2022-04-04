package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Zioyi/zocker/container"

	log "github.com/sirupsen/logrus"

	_ "github.com/Zioyi/zocker/nsenter"
)

const ENV_EXEC_PID = "zocker_pid"
const ENV_EXEC_CMD = "zocker_cmd"

func ExecContainer(containerName string, commandArray []string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Exec container getContainerPidByName %s error %v", containerName, err)
	}
	cmdStr := strings.Join(commandArray, " ")
	log.Infof("container pid %s, comand %s", pid, cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = os.Setenv(ENV_EXEC_PID, pid)
	if err != nil {
		log.Errorf("Exec container setenv %s error %v", ENV_EXEC_PID, err)
	}
	err = os.Setenv(ENV_EXEC_CMD, cmdStr)
	if err != nil {
		log.Errorf("Exec container setnv %s error %v", ENV_EXEC_CMD, err)
	}

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerName, err)
	}

}

func getContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err = json.Unmarshal(content, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}
