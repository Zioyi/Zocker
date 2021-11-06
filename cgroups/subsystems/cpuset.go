package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSetSubSystem struct{}

func (c *CpuSetSubSystem) Name() string {
	return "cpuset"
}

func (c *CpuSetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}
	return nil
}

func (c *CpuSetSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (c *CpuSetSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("set group proc fail %v", err)
	}
	return nil
}
