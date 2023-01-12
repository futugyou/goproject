package container

import (
	"encoding/json"
	"fmt"
	"golangproject/container/common"
	"io/ioutil"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func getContainerPidByName(containerName string) (string, error) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}

func getEnvsByPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %v", path, err)
	}
	return strings.Split(string(contentBytes), "\u0000"), nil
}

func getContainerInfo(containerName string) (*ContainerInfo, error) {
	filePath := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerInfoFileName)
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("read file, path: %s, err: %v", filePath, err)
		return nil, err
	}
	info := &ContainerInfo{}
	err = json.Unmarshal(bs, info)
	return info, err
}
