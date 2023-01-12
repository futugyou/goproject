package common

const (
	RootPath   = "/root/"
	MntPath    = "/root/mnt/"
	WriteLayer = "writeLayer"

	DefaultContainerInfoPath = "/var/run/go-docker/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"

	Running = "running"
	Stop    = "stopped"
	Exit    = "exited"

	EnvExecPid = "mydocker_pid"
	EnvExecCmd = "mydocker_cmd"

	IpamDefaultAllocatorPath = "/var/run/go-docker/network/ipam/subnet.json"

	DefaultNetworkPath = "/var/run/go-docker/network/network/"
)
