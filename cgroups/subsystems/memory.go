package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

func (m *MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {

	if res.MemoryLimit == "" {
		return nil
	}

	subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("set group proc fail %v", err)
	}
	return nil
}
