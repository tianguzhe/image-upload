## bili 上传图片

### 需要自己补全参数

```go
package safety

// bili
var Csrf = ""   // cookie 里的 bili_jct 字段

var Cookie = ""
```

### 使用方法

```shell

## 单文件上传
go run main.go -f [path]

## 目录上传 (指定jpg png)
go run main.go -d [dir_path]
```