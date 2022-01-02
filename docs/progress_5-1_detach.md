# 实现容器的后台运行
操作系统进程管理规则中，父进程fork出子进程后，当父进程退出时，还在运行着的子进程就变成了孤儿进程。
为了避免这些孤儿进程退出时无法释放所占用的资源而僵死，进程号为1的进程 init 就会接收这这些孤儿进程。
> 僵尸进程又是什么？
> 
> 当一个进程完成它的工作终止之后，它的父进程需要调用`wait()`或`waitpid()`获取子进程的状态信息，如果没有，这些子进程的的进程描述符仍然保存在系统中。这种进程就成为僵尸进程。

## 实现思路
解析命令参数，如果`-ti`参数为 false，说明希望这个容器在后台运行。此时我们只要在向子进程发出`zocker init ...`命令后，不去 wait 子进程，直接退出就好。

子进程会被视为一个孤儿进程，有1号进程接管。

需要注意的是：当我们以后台运行容器（子进程）时，如果在子进程中执行了`NewWorkSpace()`函数进程目录挂载，父进程在退出时执行`DeleteWorkSpace()`删除目录会报错：`Device or resource busy`。
这是因为子进程还在运行而占用着挂载的目录。

所以我们需要增加一个判断逻辑，当以后台方式运行时，不要去执行`NewWorkSpace()`。

## 运行效果
我们以 top 命令为例代表后台持续运行的一个子进程
```shell
root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ./zocker run -ti ls
INFO[0000]/vagrant_data/zocker/run.go:41 main.sendInitCommand() command all is ls
INFO[0000]/vagrant_data/zocker/main_command.go:17 main.glob..func1() init come on
INFO[0000]/vagrant_data/zocker/container/init.go:60 github.com/Zioyi/zocker/container.setUpMount() Current locaion is /root/mnt
INFO[0000]/vagrant_data/zocker/container/init.go:32 github.com/Zioyi/zocker/container.RunContainerInitProcess() Find path /bin/ls
bin   dev   etc   home  proc  root  sys   tmp   usr   var

root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ./zocker run -d top
INFO[0000]/vagrant_data/zocker/run.go:41 main.sendInitCommand() command all is top

root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
...
root      2917  2077  0 15:32 pts/0    00:00:00 sudo su
root      2918  2917  0 15:32 pts/0    00:00:00 su
root      2919  2918  0 15:32 pts/0    00:00:00 bash
root      2950     1  0 15:32 pts/0    00:00:00 top
root      2954  2919  0 15:33 pts/0    00:00:00 ps -ef
```
可以看到倒数第二行的进程2950就是我们在后台运行的"容器"，它的父进程是1

我们也可以通过`sleep 5`来模拟一个执行一会儿就结束的进程
```shell
root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ./zocker run -d sleep 5
INFO[0000]/vagrant_data/zocker/run.go:41 main.sendInitCommand() command all is sleep 5

root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
...
root      2919  2918  0 15:32 pts/0    00:00:00 bash
root      2950     1  0 15:32 pts/0    00:00:00 top
root      2977     1  0 15:41 pts/0    00:00:00 sleep 5
root      2981  2919  0 15:41 pts/0    00:00:00 ps -ef

# 等待5秒后再去检查，发现子进程消失
root@vagrant-ubuntu-trusty-64:/vagrant_data/zocker# ps -ef
UID        PID  PPID  C STIME TTY          TIME CMD
...
root      2919  2918  0 15:32 pts/0    00:00:00 bash
root      2981  2919  0 15:41 pts/0    00:00:00 ps -ef
```