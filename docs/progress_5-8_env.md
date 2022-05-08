# 实现容器指定环境变量运行
本次实现在容器运行时传入执行环境变量

## 实现思路
1. 修改容器启动命令`run`，增加`-e`选项接收用户指定的环境变量。需要支持传入多个环境变量，这里使用`cli.StringSliceFlag`接收。
然后再启动子进程的函数`NewParentProcess`中，将环境变量信息复值到`cmd.Env`中，这样在子进程启动时就会有我们设置的环境变量了
   
2. 但是如果我们以后台方式启动容器，在通过`exec`命令进入的容器时，会发现用户定义的环境变量没了。这是因为`exec`命令其实是`zocker`发起的另外一个进程，
这个进程的父进程其实是宿主机的，并不是容器的。所以我们需要将容器的环境变量从`/proc/{ContainerPID}/environ`读取出来， 
使用和上面一样的方式将环境环境变量传给新的子进程
   
## 运行效果
```shell
# 以 --ti 方式运行
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker run -ti --name bird -e bird=123 -e luck=bird busybox sh
INFO[0000]/vagrant/data/zocker/main_command.go:86 main.glob..func2() envVariables is %v[bird=123 luck=bird]       
INFO[0000]/vagrant/data/zocker/container/container_process.go:152 github.com/Zioyi/zocker/container.CreateMountPoint() mount -t aufs -o dirs=/root/writeLayer/bird:/root/busybox none /root/mnt/bird 
INFO[0000]/vagrant/data/zocker/run.go:57 main.sendInitCommand() command all is sh                            
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is bird              
INFO[0000]/vagrant/data/zocker/main_command.go:18 main.glob..func1() init come on                                 
INFO[0000]/vagrant/data/zocker/container/init.go:60 github.com/Zioyi/zocker/container.setUpMount() Current locaion is /root/mnt/bird            
INFO[0000]/vagrant/data/zocker/container/init.go:32 github.com/Zioyi/zocker/container.RunContainerInitProcess() Find path /bin/sh
/ # env | grep bird
luck=bird
SUDO_COMMAND=./zocker run -ti --name bird -e bird=123 -e luck=bird busybox sh
bird=123


# 以后台方式运行，然后通过 exec 进入
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker run -d --name bird -e bird=123 -e luck=bird busybox top
INFO[0000]/vagrant/data/zocker/main_command.go:86 main.glob..func2() envVariables is %v[bird=123 luck=bird]       
INFO[0000]/vagrant/data/zocker/container/container_process.go:152 github.com/Zioyi/zocker/container.CreateMountPoint() mount -t aufs -o dirs=/root/writeLayer/bird:/root/busybox none /root/mnt/bird 
INFO[0000]/vagrant/data/zocker/run.go:57 main.sendInitCommand() command all is top                           
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is bird              
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker ps
ID           NAME         PID         STATUS      COMMAND     CREATED
0900169580   bird         14738       running     top         2022-05-08 15:50:18
9315838849   container2   1997        running     top         2022-05-04 20:59:41
0908448955   container3   2371        running     top         2022-05-04 21:44:25
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker exec bird sh
/ # env | grep bird
luck=bird
SUDO_COMMAND=./zocker run -d --name bird -e bird=123 -e luck=bird busybox top
bird=123
```