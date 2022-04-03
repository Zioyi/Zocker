# 实现查看容器日志
查看容器日志也是一个重要的功能。当我们将容器以后台方式运行后，需要查看运行的运行情况时，日志是一个重要的方式。

## 实现思路
一般情况，容器/进程运行的日志会打印到标准输出上，我们只需要将标准输出的内容保存起来，并提供访问查看的能力即可。

我们可以将容器的标准输出重定向到指定的位置：`/var/run/zocker/{$containerName}/container.log`。
核心代码如下：
```go
cmd := exec.Command("proc/self/exe", "init")
file := os.Open("/var/run/zocker/%s/container.log", containerName)
// 这里将子进程的标准输出定向到了我们创建的日志文件中
cmd.Stdout = file
```

然后实现`zocker logs $containerName`命令，通过传入的容器名，可以构造出日志文件的路径，读取文件内容并打印到终端即可。
在此基础上，我们又实现了`-f`命令标志，如果执行`zocker logs $containerName -f`命令，将会持续刷新日志文件中的内容。


## 运行效果
我们以后台方式运行了容器`seven`
```shell
vagrant@vagrant-ubuntu-trusty-64:/vagrant_data/zocker$ sudo ./zocker run -d -name seven top
INFO[0000]/vagrant_data/zocker/run.go:59 main.sendInitCommand() command all is top
INFO[0000]/vagrant_data/zocker/run.go:46 main.Run() [zcoker] container name is seven
```
然后查看日志，可以发现内容是一致变化的
```shell
vagrant@vagrant-ubuntu-trusty-64:/vagrant_data/zocker$ sudo ./zocker logs seven -f
Mem: 335388K used, 166188K free, 388K shrd, 15084K buff, 172648K cached
CPU:  0.2% usr  0.6% sys  0.0% nic 99.1% idle  0.0% io  0.0% irq  0.0% sirq
Load average: 0.00 0.01 0.01 1/109 4
  PID  PPID USER     STAT   VSZ %VSZ CPU %CPU COMMAND
---
Mem: 335524K used, 166052K free, 392K shrd, 15104K buff, 172652K cached
CPU:  0.2% usr  0.4% sys  0.0% nic 99.3% idle  0.0% io  0.0% irq  0.0% sirq
Load average: 0.00 0.01 0.01 1/109 4
  PID  PPID USER     STAT   VSZ %VSZ CPU %CPU COMMAND
    1     0 root     R     1312  0.2   0  0.0 top
```

## TODO
这里埋了一个坑，我们对容器的日志文件没有进行回收，所以`/var/run/zocker/`文件大小会随着运行容器的增长而出现变动
之后要考虑对删除的容器做回收