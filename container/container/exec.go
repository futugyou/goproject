package container

import (
	"fmt"
	"golangproject/container/common"
	_ "golangproject/container/nsenter"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func ExecContainer(containerName string, commandArray []string) error {
	// 根据传过来的容器名获取宿主机对应的pid
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		return err
	}

	// 把命令以空格为分隔符拼接成一个字符串，便于传递
	cmdStr := strings.Join(commandArray, " ")
	logrus.Infof("container pid %s", pid)
	logrus.Infof("command %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := os.Setenv(common.EnvExecPid, pid); err != nil {
		return fmt.Errorf("setenv %s err: %v", common.EnvExecPid, err)
	}
	if err := os.Setenv(common.EnvExecCmd, cmdStr); err != nil {
		return fmt.Errorf("setenv %s err: %v", common.EnvExecCmd, err)
	}

	envs, err := getEnvsByPid(pid)
	if err != nil {
		return err
	}
	cmd.Env = append(os.Environ(), envs...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec container %s err: %v", containerName, err)
	}
	return nil
}
