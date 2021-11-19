# 通过管道通信并识别环境变量

之前的实现是通过`/proc/self/exe init args`的方式来向子进程传递参数，执行想要命令。这样有一个缺点是不能传递特殊字符
比如：
```shell
$ sudo ./zocker run -it -m stress --vm-bytes 200m --vm-keep -m 1
```
这个命令尾部的`-m 1`参数无法被传递到子进程。

解决这个问题，可以通过**进程间管道通信**的方式





