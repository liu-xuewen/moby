## mac下直接make，但是要先运行docker

make build

## linux下先安装docker

#### 1. 失败的源码编译： 没有用，docker可能用的是vgo包管理工具， 直接make build
```shell script

# 一定要放在$GOPATH下面
CGO_ENABLED=0 go build -v -mod=vendor ./cmd/dockerd
```


## dev cloud build 成功
#### 1. dev cloud build失败处理

```shell script
# dockerfile sh permission denied
# chmod -R 777 hack/dockerfile/install/
 chmod -R 777 hack/dockerfile/install/
 chmod  777 contrib/*.sh

```

#### 2. make build

#### 3 cd bundles/binary-daemon
```shell script
chmod +x dockerd
# -D 参数是debug
./dockerd -D 


chmod +x containerd
./containerd
```

#### 4. ./dockerd -D启动后一直请求 
原因是之前部署了k8s,执行k8s重置命令即可
kubeadm reset
```shell script
DEBU[2020-08-08T19:00:40.746350547+08:00] Calling GET /version      
# /containers/json?all=1&filters={"label":{"io.kubernetes.docker.type=podsandbox":true}}&limit=0                    
DEBU[2020-08-08T19:00:41.728720627+08:00] Calling GET /containers/json?all=1&filters=%7B%22label%22%3A%7B%22io.kubernetes.docker.type%3Dpodsandbox%22%3Atrue%7D%7D&limit=0 
DEBU[2020-08-08T19:00:41.731006091+08:00] Calling GET /containers/json?all=1&filters=%7B%22label%22%3A%7B%22io.kubernetes.docker.type%3Dcontainer%22%3At

```

#### 5.查看docker默认文件系统类型overlay2
```shell script
[root@VM_54_88_centos ~]# docker info | grep "Storage Driver"
 Storage Driver: overlay2

[root@VM_54_88_centos ~]# docker info | grep "Network"
  Network: bridge host ipvlan macvlan null overlay
```

#### 6. docker.service启动位置及内容
```shell script
cat /usr/lib/systemd/system/docker.service

# 源代码文件位置
docker/contrib/init/systemd
```

#### 7. 前提条件是Requires=docker.socket
```shell script
# 直接运行 /usr/bin/dockerd -H fd:// 会报错
#
#[root@VM_54_88_centos /run]# /usr/bin/dockerd -H fd://
 #INFO[2020-08-08T20:03:51.720690980+08:00] Starting up                                  
 #failed to load listeners: no sockets found via socket activation: make sure the service was started by systemd
#
# systemd 配置里面 Requires=docker.socket
systemctl start docker
# 启动的docker.sock 是在/run/docker.sock

## 直接./dockerd -D
# 启动的docker.sock 是在/var/run/docker.sock
```

#### 8. // netstat |grep docker 发现两种方式都是unix 监听
