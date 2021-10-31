package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Split(text, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	_, err := os.Stat(path.Join(cgroupRoot, cgroupPath))
	if err == nil {
		return path.Join(cgroupRoot, cgroupPath), nil
	}

	if autoCreate && os.IsNotExist(err) {
		err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755)
		if err != nil {
			return "", fmt.Errorf("error create cgroup %v", err)
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	}

	return "", fmt.Errorf("cgroup path error %v", err)
}
