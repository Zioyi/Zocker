package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(cmdArray []string, args []string) error {
	// not mount, it works
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Warnf("Not found path, %v", cmdArray[0])
		return err
	}

	logrus.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
