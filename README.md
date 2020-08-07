## mac下直接make，但是要先运行docker

make build

## linux下尝试go build 
```
// 一定要放在$GOPATH下面
CGO_ENABLED=0 go build -v -mod=vendor ./cmd/dockerd
```
没有用，docker可能用的是vgo包管理工具， 直接make build


## dev cloud build失败
1. chmod -R 777 hack/dockerfile/install/
2. chmod  777 contrib/*.sh
```
dockerfile sh permission denied
// chmod -R 777 hack/dockerfile/install/


```

