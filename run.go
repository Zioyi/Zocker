package main

import (
	"os"
	"strings"

	"github.com/Zioyi/zocker/cgroups"
	"github.com/Zioyi/zocker/cgroups/subsystems"
	"github.com/Zioyi/zocker/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig) {
	command := strings.Join(cmdArray, " ")
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 启动时限制
	cgroupManager := cgroups.NewCgroupManger("zocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	parent.Wait()
	os.Exit(-1)
}
