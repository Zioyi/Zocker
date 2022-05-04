package main

import (
	"fmt"
	"os/exec"

	"github.com/Zioyi/zocker/container"

	log "github.com/sirupsen/logrus"
)

func commitContainer(containerName string, imageName string) {
	mntURL := fmt.Sprintf(container.MntUrl, containerName)
	mntURL += "/"
	imageTar := container.RootUrl + "/" + imageName + ".tar"
	fmt.Printf("%s\n", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error %v", mntURL, err)
	}
}
