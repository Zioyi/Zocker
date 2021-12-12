package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {
	// not mount, it works
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Warnf("Not found path, %v", cmdArray[0])
		return err
	}

	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	// uintptr(3) 就是指 index 为 3 的文件描述符，也就是传递进来的管道一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

/*
init 挂载点
*/
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get currenct location error %v", err)
		return
	}
	log.Infof("Current locaion is %s", pwd)
	err = pivotRoot(pwd)
	if err != nil {
		log.Errorf("pivotRoot error %v", err)
		return
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("mount proc error %v", err)
		return
	}
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		log.Errorf("mount tmpfs error %v", err)
		return
	}
}

func pivotRoot(root string) error {
	/*
		systemd 加入linux之后，mount namespace 就变成 shared by default，所以你必须声明你要这个新的mount namespace独立。
	*/
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount new namespace error:%v", err)
	}
	/*
		为了使当前的root系统和新的root系统不在同一个mount namespace中，我们把当新的root系统重新mount一下，
		bind mount是把相同内容换了一个挂载点重新挂载
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0700); err != nil {
		return err
	}
	// pivot_root 到新的rootfs，现在老的old_root是挂载到rootfs/.pivot_root
	//
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// unmount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	return os.Remove(pivotDir)
}
