package container

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"

	"golangproject/container/common"
)

// 停止容器，修改容器状态
func StopContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info, err: %v", err)
		return
	}
	if info.Pid != "" {
		pid, _ := strconv.Atoi(info.Pid)
		// 杀死进程
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			logrus.Errorf("stop container, pid: %d, err: %v", pid, err)
			return
		}
		// 修改容器状态
		info.Status = common.Stop
		info.Pid = ""
		bs, _ := json.Marshal(info)
		fileName := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerInfoFileName)
		err := os.WriteFile(fileName, bs, 0622)
		if err != nil {
			logrus.Errorf("write container config.json, err: %v", err)
		}
	}
}
