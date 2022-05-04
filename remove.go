package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Zioyi/zocker/container"
)

func RemoveContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("[RemoveContainer] get container info by name [%s] failed, err = %v", containerName, err)
		return
	}

	if containerInfo.Status != container.Stop {
		log.Errorf("[RemoveContainer] container status [%s] is not stopped", containerInfo.Status)
		return
	}

	containerDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(containerDir); err != nil {
		log.Errorf("[RemoveContainer] remove container dir [%s] failed, err = %v", containerDir, err)
		return
	}
	container.DeleteWorkSpace(containerInfo.Volume, containerName)
}
