## mac下直接make，但是要先运行docker

？make build

## linux下尝试go build 
```
// 一定要放在$GOPATH下面
CGO_ENABLED=0 go build -v -mod=vendor ./cmd/dockerd
```
没有用，docker可能用的是vgo包管理工具，
