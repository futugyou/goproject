package container

import (
	"encoding/json"
	"fmt"
	"golangproject/container/common"
	"io/ioutil"
	"strings"
)

func getContainerPidByName(containerName string) (string, error) {
	dirUrl := fmt.Sprintf(common.DefaultContainerInfoPath, containerName)
	configFilePath := dirUrl + common.ContainerInfoFileName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("read file %s err: %v", configFilePath, err)
	}

	var containerInfo ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", fmt.Errorf("json ummarshal err: %v", err)
	}
	return containerInfo.Pid, nil
}

func getContainerInfoByName(containerName string) (*ContainerInfo, error) {
	dirUrl := fmt.Sprintf(common.DefaultContainerInfoPath, containerName)
	configFilePath := dirUrl + common.ContainerInfoFileName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %v", configFilePath, err)
	}

	var containerInfo ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return nil, fmt.Errorf("json ummarshal err: %v", err)
	}
	return &containerInfo, nil
}

func getEnvsByPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %v", path, err)
	}
	return strings.Split(string(contentBytes), "\u0000"), nil
}
