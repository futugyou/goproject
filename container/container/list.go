package container

import (
	"fmt"
	"golangproject/container/common"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
)

func ListContainerInfo() {
	files, err := ioutil.ReadDir(common.DefaultContainerInfoPath)
	if err != nil {
		logrus.Errorf("read info dir, err: %v", err)
	}

	var infos []*ContainerInfo
	for _, file := range files {
		info, err := getContainerInfo(file.Name())
		if err != nil {
			logrus.Errorf("get container info, name: %s, err: %v", file.Name(), err)
			continue
		}
		infos = append(infos, info)
	}

	// 打印
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 2, ' ', 0)
	_, _ = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Id, info.Name, info.Pid, info.Status, info.Command, info.CreateTime)
	}

	// 刷新标准输出流缓存区，将容器列表打印出来
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush info, err: %v", err)
	}
}
