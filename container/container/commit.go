package container

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func CommitContainer(containerName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Run get pwd err: %v", err)
	}
	mntUrl := pwd + "/mnt/" + containerName
	imageTar := pwd + "/" + containerName + ".tar"
	logrus.Infof("commit file path: %s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntUrl, ".").CombinedOutput(); err != nil {
		return fmt.Errorf("tar folder err: %s, %v", mntUrl, err)
	}
	logrus.Infof("end commit file: %s", imageTar)
	return nil
}
