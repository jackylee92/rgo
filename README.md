# RGO

## 简介：
锐锢-rgo

## 环境配置

* 下载安装Go，版本1.15以上，建议1.16【必须】

* 开启GO111MODULE，更换GOPROXY，网上有很多修改的方法。【必须】

  通过go env查看配置

  ````
  GO111MODULE="auto" // 或者 on
  GOPROXY="https://goproxy.cn,direct"
  ````

## 运行demo项目

* 进入你存放Go项目文件夹中，下载rgo项目【注意下载下来后的文件名一定为rgo.com】

  ```
  git clone http://git.ruigushop.com/golang/rgo-framework.git rgo.com
  ```

* 添加demo数据库，数据库文件:demo/storage/database/demo_test.sql

* 检查配置文件``config/dev.yaml`` 配置说明见注释

* 进入app/demo目录，运行

  ```shell
  cd rgo.com/app/demo_cron
  go run cmd/main.go -config=config/dev.yaml
  ```

  初次会下载所有依赖，如果出现github.com相关依赖timeout，多尝试几次，或者有条件可以开代理，再尝试。

* 启动成功后会出现 ``启动成功：启动项【bootstrap-init】【port:8228】:成功``提示。

* 测试

  使用GET方式请求http://127.0.0.1:8228/detail?id=1

  出参，测试结果数据取决于数据库数据。

  ```json
  {
      "code": 200,
      "data": {
          "config": "ljd-update",
          "name": "wupoming"
      },
      "message": "请求成功"
  }
  ```

## 创建项目

__1～3步骤为管理员执行，创建项目__

__4~9步骤为开发人员执行，编写代码__

1. 新建一个子项目仓库

2. clone 脚手架，克隆下来名称为rgo.com，【注意下载下来后的文件名一定为rgo.com】

````
git clone http://git.ruigushop.com/golang/rgo-framework.git rgo.com
````

3. 添加子项目到rgo中

在rgo.com目录下执行，因为子项目地址使用的是相对路径为app/xxxx下，如果执行位置错误，会导致子项目存放地址错误。

```
git submodule add http://git.ruigushop.com/golang/coupon-admin.git app/coupon-admin
```

4. 编写子项目代码，可以先将demo复制中案例复制过来，再次修改

   不同人员开发，需要先执行第2步，然后进入项目将子项目仓库添加到子项目文件夹中。clone如果提示子项目文件夹已存在，则可以调整clone命令，将项目的内容直接clone到子项目文件夹中

   ```
   git clone 子项目地址 ./
   ```

5. 子项目开启mod，

进入子项目，执行一下命令

````
go mod init
````

6. 子项目mod添加引入脚手架

修改子项目下的go.mod，在最后面添加如下

````
replace (
    rgo.com/bootstrap => ../../bootstrap
    rgo.com/util => ../../util
)
````

7. 整理下载mod依赖

````
go mod tidy
go mod download
````

8. 在子项目中执行

````
go run cmd/main.go -config=config/dev.yaml
````

结束

## 编译
### Mac下编译Linux, Windows

````
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
 
# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
如: CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o helloworld-windows helloworld.go
````
