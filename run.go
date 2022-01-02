package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Zioyi/zocker/cgroups"
	"github.com/Zioyi/zocker/cgroups/subsystems"
	"github.com/Zioyi/zocker/container"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 启动时限制
	cgroupManager := cgroups.NewCgroupManger("zocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)
	if tty {
		_ = parent.Wait()
		mntURL := "/root/mnt/"
		rootURL := "/root/"
		container.DeleteWorkSpace(rootURL, mntURL, volume)
	}
	os.Exit(0)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
