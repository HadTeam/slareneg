> 注意：本部分的各个部分按照功能来划分

后端分为以下几个部分
- Websocket Picker
- Game Judge
- Data Handler
- Authoritarian Handler

## Websocket Picker
外负责与 Websocket Client 通讯，内与 Etcd 交互

### Websocket Server
通过 Gorouting 与多个客户端建立连接。

### Etcd Operator
#### 游戏指令
覆盖此玩家当前回合的游戏操作。

## Game Judge
### 回合机制
通过一个计时器，在特定一段时间内结算旧回合、发起新回合。
### 多玩家操作冲突判定
在 Etcd 内存储一个 timestamp，表示 Etcd Operator 收到此游戏指令的时刻，按照 timestamp 升序排序。

## Data Handler
### 游戏结束
抽取 Etcd 中的操作信息，转存到 Mysql。

### 更新评级
定期抽取 Mysql 中的信息，计算排名。

## Authoritarian Handler
### 用户操作
包括以下
- 注册
- 登录
- 修改密码
- 恢复密码
- 绑定邮箱
- 修改邮箱

### 用户验证
被调用的地方预期有
- Websocket Server
- Web API
