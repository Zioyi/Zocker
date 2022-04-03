package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Zioyi/zocker/container"
	log "github.com/sirupsen/logrus"
)

func LogContainer(containerName string, follow bool) {
	containerLogDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	containerLogFilePath := containerLogDir + container.ContainerLogFile

	var old []byte
	for {
		content, err := ioutil.ReadFile(containerLogFilePath)
		if err != nil {
			log.Errorf("read file %s error %v", containerLogFilePath, err)
			return
		}

		if bytes.Compare(content, old) != 0 {
			_, err := fmt.Fprint(os.Stdout, string(content))
			if err != nil {
				log.Errorf("fprint to stdout error %v", err)
				return
			}
		}

		if !follow {
			break
		}
		old = content

		time.Sleep(time.Second)
	}

}
