# pbff共识算法

## 介绍

基于golang的pbft共识算法，主要用于演示PBFT共识算法的核心逻辑

1. 主节点定时出块（每10秒钟）
2. 共识流程：
    * propose - 主节点提议
    * vote - 从节点投票
    * preCommit - 所有节点预提交
    * commit - 提交，流程结束，进入下一个区块

## 编译
```
make clean
make
```

## 启主节点
```
./pbft -c=./config/config.toml
```

## 启动从节点
```
./pbft -c=./config/config1.toml
./pbft -c=./config/config2.toml
./pbft -c=./config/config3.toml
```

## 输出
```
{"level":"info","ts":1600916693.3426318,"caller":"log/zap.go:20","msg":"block commit","height":1,"data":"This is a block data, height is 1."}
{"level":"info","ts":1600916703.341808,"caller":"log/zap.go:20","msg":"block commit","height":2,"data":"This is a block data, height is 2."}
{"level":"info","ts":1600916713.339083,"caller":"log/zap.go:20","msg":"block commit","height":3,"data":"This is a block data, height is 3."}
{"level":"info","ts":1600916723.3332858,"caller":"log/zap.go:20","msg":"block commit","height":4,"data":"This is a block data, height is 4."}
{"level":"info","ts":1600916733.341336,"caller":"log/zap.go:20","msg":"block commit","height":5,"data":"This is a block data, height is 5."}
...
```