package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
}

func (c *CpuSubSystem) Name() string {
	return "cpu"
}

func (c *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuShare == "" {
		return nil
	}

	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu share fail %v", err)
	}

	return nil
}

func (c *CpuSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (c *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}
