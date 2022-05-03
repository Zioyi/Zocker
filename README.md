# Write a zocker like Docker
跟随 [xianlubird](https://github.com/xianlubird/mydocker) 书和代码，一步步实现一个容器引擎，以加深对容器技术的了解。

## Prerequirement
- 一个真实的 Linux 系统作为开发和运行环境，如：ubuntu。代码使用的一些系统调用基于 Linux 内核
- Go 开发环境
- 对 Docker 有使用经验并了解一些原理，这样在开发时会可以加强对比（非必须）

## Progress
- [x] 构造容器
    - [x] 构造实现 run 命令版本的容器 `tag: p3.1`
    - [x] 增加容器资源限制 `tag: p3.2`
    - [x] 增加管道及环境变量识别 `tag: p3.3`
- [x] 构造镜像
    - [x] 使用 busybox 创建容器 `tag: p4.1`
    - [x] 使用 AUFS 包装 busybox `tag: p4.2`
    - [x] 实现 volume 数据卷 `tag: p4.3`
    - [x] 实现简单镜像打包 `tag: p4.4`
- [ ] 构造容器进阶
    - [x] 实现容器的后台运行 `tag: p5.1`
    - [x] 实现查看运行中容器 `tag: p5.2`
    - [x] 实现查看容器日志 `tag: p5.3`
    - [x] 实现进入容器Namespace `tag: p5.4`
    - [x] 实现停止容器 `tag: p5.5`
    - [x] 实现删除容器 `tag: p5.6`
    - [ ] 实现通过容器制作镜像
    - [ ] 实现容器指定环境变量运行 

