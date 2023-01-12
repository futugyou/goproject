package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 本容器执行的第一个进程
// 使用mount挂载proc文件系统
// 以便后面通过`ps`等系统命令查看当前进程资源的情况
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("get user command in run container")
	}
	// 挂载
	err := setUpMount()
	if err != nil {
		logrus.Errorf("set up mount, err: %v", err)
		return err
	}

	// 在系统环境 PATH中寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("look %s path, err: %v", cmdArray[0], err)
		return err
	}

	err = syscall.Exec(path, cmdArray[0:], os.Environ())
	if err != nil {
		return err
	}
	return nil
}

func readUserCommand() []string {
	// 进程默认三个管道，从fork那边传过来的就是第四个（从0开始计数）
	// 指 index 为 3的文件描述符，
	// 也就是 cmd.ExtraFiles 中 我们传递过来的 readPipe
	pipe := os.NewFile(uintptr(3), "pipe")
	bs, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("read pipe, err: %v", err)
		return nil
	}
	msg := string(bs)
	return strings.Split(msg, " ")
}

func setUpMount() error {
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current location err: %v", err)
	}
	logrus.Infof("current location: %s", pwd)

	err = privotRoot(pwd)
	if err != nil {
		return err
	}

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("mount proc, err: %v", err)
		return err
	}
	syscall.Mount("tmpfs", "/dev", "tempfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	return nil
}

func privotRoot(root string) error {
	// 为了使当前root的老root和新root不在同一个文件系统下，我们把root重新mount一次
	// bind mount 是把相同的内容换了一个挂载点的挂载方法
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs、.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	// 判断当前目录是否已有该文件夹
	if _, err := os.Stat(pivotDir); err == nil {
		// 存在则删除
		if err := os.Remove(pivotDir); err != nil {
			return err
		}
	}
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir of pivot_root err: %v", err)
	}

	// pivot_root 到新的rootfs，老的old_root现在挂载在rootfs/.pivot_root上
	// 挂载点目前依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root err: %v", err)
	}

	// 修改当前工作目录到跟目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir root err: %v", err)
	}

	// 取消临时文件.pivot_root的挂载并删除它
	// 注意当前已经在根目录下，所以临时文件的目录也改变了
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir err: %v", err)
	}
	return os.Remove(pivotDir)
}
