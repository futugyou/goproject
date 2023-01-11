### mount hierarchy, but is hierarchy not link with any subsystem, so it can not use cgroup to limit resouce. 
```
mkdir cgroup-demo
mount -t cgroup -o none,name=cgroup-demo cgroup-demo ./cgroup-demo

ls cgroup-demo
cgroup.clone_children: if 1 use parent cpuset
cgroup.procs:  process set id
task: process id
```

### /sys/fs/cgroup, subsystem with hierarchy will be here.
```
cd /sys/fs/cgroup/memory
mkdir cgroup-demo-memory

ls cgroup-demo-memory will see many file, we can write pid to 'tasks' and modify meory.limit_in_bytes to limit meory use.
```

### read command line args
```
go get github.com/urfave/cli
```

### log
```
go get github.com/sirupsen/logrus
```

### export 
```
docker pull busybox
docker run -d busybox top -b
docker ps
docker export -o busybox.tar  'ID'
mkdir busybox && tar -xvf busybox.tar -C busybox/
```

### mount 
 mount -t aufs -o dirs=/root/writeLayer:/root/busybox none /root/mn