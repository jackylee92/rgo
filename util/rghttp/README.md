# HTTP工具

## 使用

```go
// 地址
url := "http:xxxx"
// 参数
param := `{"id": 1}`
// 实例化一个请求对象
httpClient := rghttp.Client {
    Param : param,
    Method: "POST", // 请求方式
    Header: map[string]string{ // header头
        "Content-Type": "application/json",
        },
    Url: url,
    This: this // 上下文对象this
    }

httpCode, data, err := httpClient.GetApi()
if err != nil {
  // 请求错误了
}

```
