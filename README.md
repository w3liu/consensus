# consensus
consensus base golang

## 编译
```
cd pbft
go build
```

## 运行
```
# 启动第一个节点
./pbft.exe -c=config/config.toml

# 启动第二个节点
./pbft.exe -c=config/config1.toml

# 启动第三个节点
./pbft.exe -c=config/config2.toml

# 启动第四个节点
./pbft.exe -c=config/config3.toml
```

## 效果
[!qr](./docs/images/pbft.png)
