# RGO

## 简介：
rgo

# Content
* [说明](#说明)
* [安装](#安装)
* [使用](#使用)

## 说明
## 安装
## 使用
## 环境配置

- 1.golang version >= 1.16
- 2.开启gomodule
```
go env -w GO111MODULE=on

go env -w GOPROXY=https://goproxy.cn,direct
```
- 3.配置私有仓库
```
export GOPRIVATE=git.ruigushop.com
go env -w GOINSECURE=git.ruigushop.com
go env -w GOPRIVATE=git.ruigushop.com
```

- 4.修改git配置 
```
vi ~/.gitconfig 

[url "git@git.ruigushop.com:"]
        insteadOf = https://git.ruigushop.com/
[http]
        extraheader = PRIVATE-TOKEN:xxxxxxxxxxxxxxxxxx
```
- 5.当git配置不生效时修改 netrc
```
vim ~/.netrc
machine git.ruigushop.com
    login xxxxx
    password xxxxxxxxxxxxxxxxxx
```

- 6.在项目中使用
```
go get github.com/jackylee92/rgo
```