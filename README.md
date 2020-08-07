## mac下直接make，但是要先运行docker

make build

## linux下先安装docker

#### 1. 失败的源码编译： 没有用，docker可能用的是vgo包管理工具， 直接make build
```

// 一定要放在$GOPATH下面
CGO_ENABLED=0 go build -v -mod=vendor ./cmd/dockerd
```


## dev cloud build 成功
#### 1. dev cloud build失败处理
1. chmod -R 777 hack/dockerfile/install/
2. chmod  777 contrib/*.sh
```
dockerfile sh permission denied
// chmod -R 777 hack/dockerfile/install/


```

#### 2. make build

#### 3 cd bundles/binary-daemon
```
chmod +x dockerd
./dockerd


chmod +x containerd
./containerd
```
