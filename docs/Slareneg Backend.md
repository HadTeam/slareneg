> 注意：本部分的各个部分按照功能来划分

后端分为以下几个部分
- Api Provider
  - Websocket Receiver
  - HTTP Receiver
  - Game Operator
- Game Pool
  - Game Judge
- Data Saver

共有两个储存接口
- Temperature Data Source
- Persistent Data Source

还有 Rpc 接口
- Api Provider Rpc
- Game Pool Rpc

## Data Source
### Temperature
- 获取目前地图
- 更新地图
- 获取指令列表
- 更新游戏指令
- 修改游戏状态
- 获取玩家列表
- 修改玩家状态

### Persistent
- 获取初始地图

## Api Provider

### Websocket Receiver
与 Websocket Client 通讯以接收指令，并与 Game Operator 交互以执行指令。

### HTTP Receiver
负责相应 HTTP API 请求。

### Authoritarian Handler
#### 用户操作
包括以下
- 注册
- 登录
- 修改密码
- 恢复密码
- 绑定邮箱
- 修改邮箱

## Game Pool
### 回合机制
通过一个计时器，在特定一段时间内结算旧回合、发起新回合。

### 多玩家操作冲突判定
通过 Instruction Temp 机制与按时间戳排序来实现

Instruction Temp 机制：在 Temperature Data Source 内设置一个回合数字段，按照回合数与玩家 ID 覆盖保存 Instruction，当回合结束时使回合数递增。

## Data Saver
### 游戏结束时转存数据
抽取 Temp DataSource 中的操作信息，转存到 Persistent Data Source。

### 更新评级
定期抽取 Persistent Data Source 中的信息，计算排名。