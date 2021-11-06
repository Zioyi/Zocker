# Cgroups
它提供了一种控制进程（及子进程）的资源限制、控制和统计能力。
资源包含CPU、内存、存储、网络等。

cgroup是对进程分组管理的一种机制，一个cgroup包含一组进程。

## cgroups子系统（subsystem）
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

## cgroups 层级结构（hierarchy）
`hierarchy`的功能是把一组 cgroup 组织成一个树状结构，让 Cgroup 可以实现继承
> 一个 cgroup1 限制了其下的进程（P1、P2、P3）的 CPU 使用频率，如果还想对进程P2进行内存的限制，可以在 cgroup1 下创建 cgroup2，使其继承于 cgroup1，可以限制 CPU 使用率，又可以设定内存的限制而不影响其他进程。

内核使用 cgroups 结构体来表示对某一个或某几个 cgroups 子系统的资源限制，它是通过一棵树的形式进行组织，被成为`hierarchy`.


## cgroups 与进程
hierarchy、subsystem 与cgroup进程组间的关系
hierarchy 只实现了继承关系，真正的资源限制还是要靠 subsystem
通过将 subsystem 附加到 hierarchy上，
将进程组 加入到 hierarchy下（task中），实现资源限制

![image](https://awps-assets.meituan.net/mit-x/blog-images-bundle-2015/3982f44c.png)

通过这张图可以看出：
- 一个 subsystem 只能附加到一个 hierarchy 上面
- 一个 hierarchy 可以附加多个 subsystem
- 一个进程可以作为多个 cgroup 的成员，但是这些 cgroup 必须在不同hierarchy 中。
- 一个进程 fork 出子进程时，子进程是和父进程在同一个 cgroup 中的，也可以根据需要将其移动到其他 cgroup 中。

## cgroups 文件系统
cgroups 的底层实现被 Linux 内核的 VFS（Virtual File System）进行了隐藏，给用户态暴露了统一的文件系统 API 借口。我们来体验一下这个文件系统的使用方式：

1. 首先，要创建并挂载一个hierarchy（cgroup树） 
```shell
$ mkdir cgroup-test
$ sudo mount -t cgroup -o none,name=cgroup-test cgrout-test ./cgroup-test
$ ls ./cgrpup-test
cgroup.clone_children  cgroup.sane_behavior  release_agent
cgroup.procs           notify_on_release     tasks
```
这些文件就是这个hierarchy中cgroup根节点的配置项

cgroup.clone_children 会被 cpuset 的 subsystem 读取，如果是1，子 cgroup 会继承父 cgroup 的 cpuset 的配置。

notify_on_release 和 release_agent 用于管理当最后一个进程退出时执行一些操作

tasks 标识该 cgroup 下面的进程 ID，将 cgroup 的进程成员与这个 hierarchy 关联

2.再创建两个子 hierarchy创建刚刚创建好的hierarchy上cgroup根节点中扩展出的两个子cgroup
```shell
$ cd cgroup-test
$ sudo mkdir cgroup-1
$ sudo mkdir cgroup-2
$ tree
.
├── cgroup-1
│   ├── cgroup.clone_children
│   ├── cgroup.procs
│   ├── notify_on_release
│   └── tasks
├── cgroup-2
│   ├── cgroup.clone_children
│   ├── cgroup.procs
│   ├── notify_on_release
│   └── tasks
├── cgroup.clone_children
├── cgroup.procs
├── cgroup.sane_behavior
├── notify_on_release
├── release_agent
└── tasks

2 directories, 14 files
```
可以看到，在一个 cgroup 的目录下创建文件夹时，Kernel 会把文件夹标记为这个 cgroup 的子 cgroup，它们会继承父 cgroup 的属性。

3. 向cgroup中添加和移动进程
一个进程在一个Cgroups的hierarchy中，只能在一个cgroup节点上存在，系统的所有进程都会默认在根节点上存在，可以将进程移动到其他cgroup节点，只需要将进程ID写到移动到的cgroup节点的tasks文件中即可。
```shell
# cgroup-test
$ ehco $$
3444
$ cat /proc/3444/cgroup 
13:name=cgroup-test:/
12:cpuset:/
11:rdma:/
10:devices:/user.slice
9:perf_event:/
8:net_cls,net_prio:/
7:pids:/user.slice/user-1000.slice/user@1000.service
6:memory:/user.slice/user-1000.slice/user@1000.service
...
```
可以看到当前终端的进程在根 cgroup 下，我们现在把他移动到子 cgroup 下
```shell
$ cd cgroup-1
$ sudo sh -c "echo $$ >> tasks"
$ cat /proc/3444/cgroup
13:name=cgroup-test:/cgroup-1
12:cpuset:/
11:rdma:/
10:devices:/user.slice
9:perf_event:/
8:net_cls,net_prio:/
7:pids:/user.slice/user-1000.slice/user@1000.service
6:memory:/user.slice/user-1000.slice/user@1000.service
...
```
可以看到终端进程所属的 cgroup 已将变成了 cgroup-1，再看一下父 cgroup 的tasks，已经没有了终端进程的 ID
```shell
$ cd cgroup-test
$ cat tasks | grep "3444"
# 返回为空
```

4. 通过 subsystem 限制 cgroup 中进程的资源。

操作系统默认已为每一个 subsystem 创建了一个默认的 hierarchy，在`sys/fs/cgroup/`目录下
```shell
$ ls /sys/fs/cgroup
blkio    cpu,cpuacct  freezer  net_cls           perf_event  systemd
cpu      cpuset       hugetlb  net_cls,net_prio  pids        unified
cpuacct  devices      memory   net_prio          rdma
```
可以看到内存子系统的 hierarchy 也在其中创建一个子cgroup
``` shell
$ cd /sys/fs/cgroup/memory
$ sudo mkdir test-limit-memory && cd test-limit-memorysudo
# 设置最大内存使用为 100MB
$ sudo sh -c "echo "100m" > memory.limit_in_bytes"sudo sh -c "echo $$ > tasks"
sudo sh -c "echo $$ > tasks"
$ sudo sh -c "echo $$ > tasks"
# 运行占用内存200MB 的 stress 经常
$ stress --vm-bytes 200m --vm-keep -m 1
```
可以对比运行前后的内存剩余量，大概只减少了100MB
```shell
# 运行前
$ top
top - 12:04:12 up  6:45,  1 user,  load average: 1.87, 1.29, 1.06
任务: 348 total,   1 running, 346 sleeping,   0 stopped,   1 zombie
%Cpu(s):  1.3 us,  0.9 sy,  0.0 ni, 97.7 id,  0.0 wa,  0.0 hi,  0.1 si,  0.0 st
MiB Mem :   5973.4 total,    210.8 free,   2820.9 used,   2941.8 buff/cache
MiB Swap:    923.3 total,    921.9 free,      1.3 used.   2746.3 avail Mem 

# 运行后
$ top
top - 12:04:57 up  6:45,  1 user,  load average: 2.25, 1.44, 1.12
任务: 351 total,   3 running, 347 sleeping,   0 stopped,   1 zombie
%Cpu(s): 34.3 us, 32.8 sy,  0.0 ni, 21.1 id,  4.9 wa,  0.0 hi,  6.9 si,  0.0 st
MiB Mem :   5973.4 total,    118.6 free,   2956.7 used,   2898.1 buff/cache
MiB Swap:    923.3 total,    817.7 free,    105.5 used.   2604.5 avail Mem 
```
说明 cgroup 的限制生效了

## docker 中是怎样进行 cgroup 限制的
首先运行一个被限制内存的容器
```shell
$ sudo docker pull redis:4
$ sudo docker run -tid -m 100m redis:4
d79f22eb11d22c56a90f88e0aeb3cfda7cbe9639e2ab0e8532003a695e375e8d
```
查看原来的内存子系统绑定的cgroup，会看到里面多了子cgroup, `docker` 
```shell
$ ls /sys/fs/cgroup/memory
... docker
...
$ ls /sys/fs/cgroup/memory/docker
cgroup.clone_children                                             memory.max_usage_in_bytes
cgroup.event_control                                              memory.memsw.failcnt
cgroup.procs                                                      memory.memsw.limit_in_bytes
d79f22eb11d22c56a90f88e0aeb3cfda7cbe9639e2ab0e8532003a695e375e8d  memory.memsw.max_usage_in_bytes
memory.failcnt                                                    memory.memsw.usage_in_bytes
memory.force_empty                                                memory.move_charge_at_immigrate
memory.kmem.failcnt                                               memory.numa_stat
memory.kmem.limit_in_bytes                                        memory.oom_control
memory.kmem.max_usage_in_bytes                                    memory.pressure_level
memory.kmem.slabinfo                                              memory.soft_limit_in_bytes
memory.kmem.tcp.failcnt                                           memory.stat
memory.kmem.tcp.limit_in_bytes                                    memory.swappiness
memory.kmem.tcp.max_usage_in_bytes                                memory.usage_in_bytes
memory.kmem.tcp.usage_in_bytes                                    memory.use_hierarchy
memory.kmem.usage_in_bytes                                        notify_on_release
memory.limit_in_bytes                                             tasks
```
可以看到`docker`cgroup里面的`d79f22eb11d22c56a90f88e0aeb3cfda7cbe9639e2ab0e8532003a695e375e8d`cgroup 正好是我们
刚才创建的容器 ID，那么看一下里面吧
```shell
$ cd /sys/fs/cgroup/memory/docker/d79f22eb11d22c56a90f88e0aeb3cfda7cbe9639e2ab0e8532003a695e375e8d
$ cat memory.limit_in_bytes
104857600cat
# 正好是100MB
```

## 总结
讲述了 cgroups 的原理，它是通过三个概念（cgroup、subsystem、hierarchy）进行组织和关联的，可以理解为
3层结构，将进程关联在 cgroup 中，然后把 cgroup 与 hierarchy 关联，subsystem 再与 hierarchy 关联，从而在限制进程资源的基础上达到一定
的复用能力。

讲述了 docker 的具体实现方式，在使用 docker 时，也能从心中了然它时怎么做到对容器使用资源的限制的。