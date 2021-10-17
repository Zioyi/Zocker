# 构造实现 run 命令版本的容器

## 基础知识
`what happen when exec "docker run -ti iamge-foo /bin/bash" in shell?`
容器在被 Docker 启动后是以进程的形式运行在 Docker 所在的宿主机系统上的，
并通过 Linux 内核提供的隔离(Namespace)和限制(Cgroup)能力，承诺不同容器间是“不可见的”，就像原来通过虚拟机来运行一样。

这里，我们将启动容器简单想象成启动一个新进程，在启动后要求：
1. 新进程被 Namespace 所限定，他只能看到他能看到
2. 以新进程视角来看，他的PID为1（即他代表这个空间里的跟进程）

```shell
root ps aux
	|
	|
	|--PID 777 docker run -it image-bar /bin/bash
	|
	|--PID 999 docker run -it image-foo /bin/bash
		 |($ ps aux)
		 |
		 |--PID 1 /bin/bash
		 |--PID 2 ps -ef
			
```

通过 Go 基础库`os/exec`提供的能力，可以在 fork 新进程时，通过`Clnoeflags`参数按需设定 Namespace 来对新进程进行隔离。

> 与 namesapce 相关的三个系统调用
>- clone() 创建新进程。提供参数来制定创建哪些 Namespace
>- unsahre() 将进程移出某个 NameSpace
>- setns() 将进程加入到 NameSpace 中


通过 Go 基础库`syscall`对系统调用`execve`的封装，可以达到覆盖当前进程的PID等信息。

## 未解决问题
在 zocker 创建好隔离的新进程后，在新进程中调用初始化函数，会挂载proc文件系统。即便是退出容器，挂载仍没有恢复，必须重启

