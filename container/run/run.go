package run

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"golangproject/container/cgroups"
	"golangproject/container/cgroups/subsystem"
	"golangproject/container/container"
	"golangproject/container/network"
)

func Run(cmdArray []string, tty, detach bool, res *subsystem.ResourceConfig, containerName, imageName, volume string,
	envs []string, nw string, portMapping []string) {
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envs)
	if parent == nil {
		logrus.Errorf("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start failed, err: %v", err)
		return
	}
	// 记录容器信息
	containerID, containerName, err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, containerName)
	if err != nil {
		logrus.Errorf("record container info, err: %v", err)
	}

	// 添加资源限制
	cgroupMananger := cgroups.NewCGroupManager("go-docker")
	// 删除资源限制
	defer cgroupMananger.Destroy()
	// 设置资源限制
	cgroupMananger.Set(res)
	// 将容器进程，加入到各个subsystem挂载对应的cgroup中
	cgroupMananger.Apply(parent.Process.Pid)

	if nw != "" {
		// 定义网络
		_ = network.Init()
		containerInfo := &container.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portMapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			logrus.Errorf("connnect network err: %v", err)
			return
		}
	}

	// 设置初始化命令
	sendInitCommand(cmdArray, writePipe)

	if !detach {
		// 等待父进程结束
		err := parent.Wait()
		if err != nil {
			logrus.Errorf("parent wait, err: %v", err)
		}
		// 删除容器工作空间
		err = container.DeleteWorkSpace(containerName, volume)
		if err != nil {
			logrus.Errorf("delete work space, err: %v", err)
		}
		// 删除容器信息
		container.DeleteContainerInfo(containerName)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
