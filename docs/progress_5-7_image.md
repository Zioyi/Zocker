# 实现通过容器制作镜像
之前在[4-2节](./progress_4-2_aufs.md)实现了在启动容器时可以使用镜像的AUFS文件系统，但没有考虑到：
1. 现在镜像`busybox`是写死在代码中，不能指定其他的镜像
2. 如果使用同一个镜像进行启动多个容器，容器之间的可写层是相互影响的

同时，我们也需要支持对不同容器进行打包镜像的功能。


## 实现思路
1.解决启动容器时无法自由指定镜像的问题

修改`run`子命令，获取命令行中传入的镜像名称`imageName`参数，在之后的挂载系统时使用`imageName`
来代替之前写死的`busybox`，比如：
```go
// 以前的代码
busyboxURL := rootURL + "busybox/"
busyboxTarUrl := rootURL + "busybox.tar"
exist, err := PathExists(busyboxURL)

// 新的代码
unTarFolderUrl := RootUrl + "/" + imageName + "/"
busyboxTarUrl := RootUrl + "/" + imageName + ".tar"
exist, err := PathExists(unTarFolderUrl)
```


2.解决启动多个容器，容器之间的可写层是相互影响的问题

每个容器都有独立的挂载点：点原先的可读可写层路径是：`root/writeLayer`，我们在其下面增加以容器名命名的子目录作为各自的可读可写层，
以此达到容器之间可写层相互隔离的目的
```shell
/root/writeLayer/
|-- container1
|   |-- root
|   |-- to1
|   `-- to1-1
|       `-- test1.txt
`-- container2
    `-- to2
```

3. 解决自由打包镜像的问题
   
修改`commit`子命令，获取命令行中传入的容器名称`containerName`参数，因为容器的挂载点路径是通过容器名
构造的，并且是独立的，所有只要构造出挂载点路径进行实现打包


## 运行效果
```shell
# 首先启动两个容器，它们都有自己的各自的卷和挂载
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# sudo ./zocker run -d --name container1 -v /root/from1:/to1 busybox top
INFO[0000]/vagrant/data/zocker/container/container_process.go:151 github.com/Zioyi/zocker/container.CreateMountPoint() mount -t aufs -o dirs=/root/writeLayer/container1:/root/busybox none /root/mnt/container1 
INFO[0000]/vagrant/data/zocker/container/container_process.go:88 github.com/Zioyi/zocker/container.NewWorkSpace() ["/root/from1" "/to1"]                       
INFO[0000]/vagrant/data/zocker/run.go:57 main.sendInitCommand() command all is top                           
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is container1        
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# sudo ./zocker run -d --name container2 -v /root/from2:/to2 busybox top
INFO[0000]/vagrant/data/zocker/container/container_process.go:151 github.com/Zioyi/zocker/container.CreateMountPoint() mount -t aufs -o dirs=/root/writeLayer/container2:/root/busybox none /root/mnt/container2 
INFO[0000]/vagrant/data/zocker/container/container_process.go:88 github.com/Zioyi/zocker/container.NewWorkSpace() ["/root/from2" "/to2"]                       
INFO[0000]/vagrant/data/zocker/run.go:57 main.sendInitCommand() command all is top                           
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is container2        
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# sudo ./zocker ps
ID           NAME         PID         STATUS      COMMAND     CREATED
4053391054   container1   1977        running     top         2022-05-04 20:59:22
9315838849   container2   1997        running     top         2022-05-04 20:59:41

# 进入container1执行一些写入操作
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# sudo ./zocker exec container1 sh
INFO[0000]/vagrant/data/zocker/exec.go:27 main.ExecContainer() container pid 1977, comand sh                
got zocker_pid=1977
got zocker_cmd=sh
setns on ipc namespace succeeded
setns on uts namespace succeeded
setns on net namespace succeeded
setns on pid namespace succeeded
setns on mnt namespace succeeded
/ # ls
bin   dev   etc   home  proc  root  sys   tmp   to1   usr   var
/ # echo "hello, container1" >> to1/test1.txt
/ # mkdir to1-1
/ # echo "hello container1,to1-1,test1" >> to1-1/test1.txt
/ # exit

# 查看一个可写层，产生了对应的内容
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# sudo tree /root/writeLayer/
/root/writeLayer/
|-- container1
|   |-- root
|   |-- to1
|   `-- to1-1
|       `-- test1.txt
`-- container2
    `-- to2

6 directories, 1 file

# 与宿主机绑定的数据卷位置也有了对应的写入
root@vagrant-ubuntu-trusty-64:~# tree from1/
from1/
`-- test1.txt

0 directories, 1 file

# 我们将container1打包为镜像，解压后发现进行当时容器的变更都在
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# ./zocker commit container1 image1
/root/image1.tar
root@vagrant-ubuntu-trusty-64:~# mkdir image1
root@vagrant-ubuntu-trusty-64:~# tar -xvf image1.tar -C image1/
root@vagrant-ubuntu-trusty-64:~# cat /root/image1/to1-1/test1.txt 
hello container1,to1-1,test1
root@vagrant-ubuntu-trusty-64:~# cat /root/image1/to1/test1.txt 
hello, container1

# 我们以这个镜像再启动一个新的容器，之前容器的写入时存在
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# ./zocker run -d --name container3 -v /root/from3:/to3 image1 top
INFO[0000]/vagrant/data/zocker/container/container_process.go:151 github.com/Zioyi/zocker/container.CreateMountPoint() mount -t aufs -o dirs=/root/writeLayer/container3:/root/image1 none /root/mnt/container3 
INFO[0000]/vagrant/data/zocker/container/container_process.go:88 github.com/Zioyi/zocker/container.NewWorkSpace() ["/root/from3" "/to3"]                       
INFO[0000]/vagrant/data/zocker/run.go:57 main.sendInitCommand() command all is top                           
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is container3        
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# ./zocker exec container3 sh
INFO[0000]/vagrant/data/zocker/exec.go:27 main.ExecContainer() container pid 2371, comand sh                
got zocker_pid=2371
got zocker_cmd=sh
setns on ipc namespace succeeded
setns on uts namespace succeeded
setns on net namespace succeeded
setns on pid namespace succeeded
setns on mnt namespace succeeded
/ # ls
bin    dev    etc    home   proc   root   sys    tmp    to1    to1-1  to3    usr    var
/ # cat to
to1-1/  to1/    to3/
/ # cat to1/test1.txt 
hello, container1
/ # cat to1-1/test1.txt 
hello container1,to1-1,test1

# 在检测一下删除容器，挂载点已经消失，说明删除容器时正常的
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# ./zocker rm container1
root@vagrant-ubuntu-trusty-64:/vagrant/data/zocker# tree /root/writeLayer/
/root/writeLayer/
`-- container2
    `-- to2

```