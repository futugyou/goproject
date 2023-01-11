package demo

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Resouceisolation() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// isolation uts,ipc,pid,mount,user,network
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
		// set container UID and GID
		UidMappings: []syscall.SysProcIDMap{
			{
				// contanier UID
				ContainerID: 1,
				// host UID
				HostID: 0,
				Size:   1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				// contanier UID
				ContainerID: 1,
				// host UID
				HostID: 0,
				Size:   1,
			},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
