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