### go-kv
key & value simple database

### feature
- client和server之间通过tcp进行通信
- database可以配置持久化
- 交互式命令行客户端

### usage
- 拷贝env.yaml.example为env.yaml并完成配置
- 启动server，可以指定配置文件，默认env.yaml
```
go build -o kv kv.go
kv -c env.yaml
```
- 启动client
```
go build -o kv_cli kv_cli.go
kv_cli -h 127.0.0.1 -p 9000
```
- 客户端命令：
```
set key value
get key
delete key
list
persistent
help
exit
```

### example
![avatar](https://github.com/markbest/go-kv/blob/master/kv_cli.gif)