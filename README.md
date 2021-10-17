# Write a zocker like Docker
跟随 [xianlubird](https://github.com/xianlubird/mydocker) 书和代码，一步步实现一个容器引擎，以加深对容器技术的了解。

## Prerequirement
- 一个真实的 Linux 系统作为开发和运行环境，如：ubuntu。代码使用的一些系统调用基于 Linux 内核
- Go 开发环境
- 对 Docker 有使用经验并了解一些原理，这样在开发时会可以加强对比（非必须）

## Progress
- [ ] 构造容器
    - [x] 构造实现 run 命令版本的容器 `tag: p3.1`
    - [ ] 增加容器资源限制
    - [ ] 增加管道及环境变量识别
- [ ] 构造镜像
    - [ ] 使用 busybox 创建容器
    - [ ] 使用 AUFS 包装 busybox
    - [ ] 实现 volume 数据卷
    - [ ] 实现简单镜像打包
- [ ] TBD

