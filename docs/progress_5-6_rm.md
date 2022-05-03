# 实现容器删除
本次实现将停止的容器进行删除，即删除容器在生命周期内产生的数据。

## 实现思路
1. 在命令行中增加`rm`子命令，使用时需要传入想要删除的容器名
2. 通过容器名查出该容器信息，检查容器状态必须为`stopped`
3. 使用系统函数`os.RemoveAll`删除该容器产生的目录

## 运行效果
```shell
# 首先运行一个容器 bird
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker run -d --name bird top
INFO[0000]/vagrant/data/zocker/run.go:59 main.sendInitCommand() command all is top                           
INFO[0000]/vagrant/data/zocker/run.go:46 main.Run() [zocker] container name is bird   

# 可以看到 bird 容器建立了一个目录来存储信息
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo tree /var/run/zocker/
/var/run/zocker/
`-- bird
    |-- config.json
    `-- container.log

1 directory, 2 files

# 我们将 bird 停止
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker stop bird
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker ps
ID           NAME        PID         STATUS      COMMAND     CREATED
4511194203   bird                    stopped     top         2022-05-03 21:55:45

# 执行删除命令，再查看`var/run/zocker`目录，bird 的目录已经不在了
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo ./zocker rm bird
vagrant@vagrant-ubuntu-trusty-64:/vagrant/data/zocker$ sudo tree /var/run/zocker/
/var/run/zocker/

0 directories, 0 files
```
