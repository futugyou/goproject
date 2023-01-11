package subsystem

import (
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

type MemorySubSystem struct {
	apply bool
}

func (*MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", cgroupPath, err)
		return err
	}
	if res.MemoryLimit != "" {
		m.apply = true
		// 设置cgroup内存限制，
		// 将这个限制写入到cgroup对应目录的 memory.limit_in_bytes文件中即可
		err := os.WriteFile(path.Join(subsystemCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemorySubSystem) Remove(cgroupPath string) error {
	subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsystemCgroupPath)
}

func (m *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if m.apply {
		subsystemCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath, false)
		if err != nil {
			return err
		}
		tasksPath := path.Join(subsystemCgroupPath, "tasks")
		err = os.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks, path: %s, pid: %d, err: %v", tasksPath, pid, err)
			return err
		}
	}
	return nil
}
