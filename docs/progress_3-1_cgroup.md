# Ggroups
它提供了一种控制进程（及子进程）的资源限制、控制和统计能力。
资源包含CPU、内存、存储、网络等。

它的实现原理是通过对各类 Linux subsystem 进行参数设定，然后将进程和这些子系统进行绑定。

Linux subsystem 有以下几种：
- blkio
- cpu
- cpuacct 统计cgroup中进程的CPU占用
- cpuset
- devices
- freezer 用户挂起和恢复从group中的进程
- memeory 控制cgroup中进程的内存占用
- net_cls
- net_prio
- ns

通过安装 cgroup 工具
```shell
$ apt-get install cgroup-tools
$ lssubsys -a
cpuset
cpu,cpuacct
blkio
memory
devices
freezer
net_cls,net_prio
perf_event
hugetlb
pids
rdma
```

