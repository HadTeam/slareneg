# Slareneg

## 开发计划

### 整体
仿照 [generals.io](https://generals.io/) 完成一个网页小游戏

- generals 核心玩法
   - 地图元素（空格，兵格，城堡，将军，山脉）
   - 可操作元素（兵的移动、占领，地图缩放和拖动等）
   - 回合机制
- 简单的控制界面（模式选择）


#### 前端
- 地图元素
- 可操作元素

#### 后端
尚未考虑：
- 地图的生成
- 地图市场
- 玩家排位机制（分数计算/排行榜规则）

##### 架构设计

后端游戏系统采用分层架构设计，主要分为两个核心层次：

1. **事件转发层 (Game)**
   - 负责事件的订阅、解析和转发
   - 消费来自消息队列的player/control事件
   - 调用GameCore的相应方法处理游戏逻辑
   - 将GameCore返回的事件发布到相应的topic (broadcast/player/control)
   - 处理游戏生命周期管理

2. **游戏逻辑层 (GameCore/BaseCore)**
   - 纯粹的游戏逻辑处理，无外部依赖
   - 提供标准化的事件返回接口
   - 处理玩家动作、回合机制、胜负判定等核心逻辑
   - 使用统一的事件类型系统 (core/event.go)

这种设计的优势：
- **职责分离**：事件处理与游戏逻辑完全分离
- **可测试性**：GameCore可以独立进行单元测试
- **可扩展性**：可以轻松替换事件系统或游戏逻辑
- **类型安全**：所有事件使用统一的类型定义

##### 大致数据流
```mermaid
sequenceDiagram
    participant Client as 客户端
    participant HTTP as 认证服务
    participant WS as WebSocket网关 
    participant Lobby as 游戏大厅
    participant Game as 游戏核心服务
    participant Queue as 消息队列
    participant Cache as 缓存
    participant DB as 数据库

    %% 1. 用户认证流程
    alt 用户认证
        Client->>HTTP: 发送注册/登录请求
        HTTP->>DB: 验证用户凭据
        DB-->>HTTP: 返回用户信息
        HTTP->>Cache: 缓存用户会话
        HTTP-->>Client: 返回JWT令牌
    end

    %% 2. 建立连接与游戏进入
    Client->>WS: 携带JWT建立WebSocket连接
    WS->>Cache: 验证JWT有效性
    Cache-->>WS: 验证成功
    WS-->>Client: 连接成功确认

    Client->>WS: 请求加入游戏
    WS->>Lobby: 转发加入请求
    Lobby->>Cache: 查询/创建游戏房间
    Lobby->>Game: 创建游戏实例
    Game->>Game: 初始化游戏状态(可能要在此过程中请求数据库等以获取地图)
    
    %% 3. 玩家加入通知
    Game->>Queue: 发布玩家加入事件
    note right of Game: 发布到 `${room.id}/broadcast`
    Queue-->>WS: 订阅者接收事件
    note left of WS: 订阅 `${room.id}/broadcast`
    
    Lobby->>WS: 返回房间信息
    WS-->>Client: 成功进入房间通知
    WS-->>Client: 广播新玩家加入

    %% 4. 核心游戏循环
    alt 游戏进行中
        Client->>WS: 发送游戏指令(如移动)
        WS->>Queue: 将指令加入队列
        note right of WS: 发布到 `${game.id}/commands`

        Game->>Queue: 消费玩家指令
        note left of Game: 订阅 `${game.id}/commands`
        note right of Game: 处理流程:<br/>1. 验证指令<br/>2. 计算状态变更
        Game->>Cache: 更新游戏快照

        alt 有效指令
            Game->>Queue: 发布游戏状态更新事件
            note right of Game: 发布到 `${room.id}/broadcast`
            Queue-->>WS: 订阅者接收事件
            WS-->>Client: 广播最新状态
        else 无效指令
            Game->>Queue: 发布错误事件给特定玩家
            note right of Game: 发布到 `${player.id}/notifications`
            Queue-->>WS: 订阅者接收事件
            WS-->>Client: 发送错误提示
        end
        
        %% 游戏定时器/回合逻辑
        loop 定时触发
            Game->>Game: 内部处理回合逻辑
            Game->>Queue: 发布新回合事件
            note right of Game: 发布到 `${room.id}/broadcast`
            Queue-->>WS: 订阅者接收事件
            WS-->>Client: 广播新回合信息
        end
    end

    %% 5. 游戏结束与结算
    alt 游戏结束
        Game->>Game: 检查结束条件
        Game->>Queue: 发布游戏结束事件
        note right of Game: 发布到 `${room.id}/broadcast`
        Game->>DB: 持久化游戏结果
        Game->>Cache: 清理游戏缓存
        Queue-->>WS: 订阅者接收事件
        WS-->>Client: 推送结算信息
    end

    %% 6. 异常处理
    alt 玩家断线
        WS->>Lobby: 通知玩家断线
        Lobby->>Game: 通知游戏实例
        Game->>Queue: 发布玩家离开事件
        note right of Game: 发布到 `${room.id}/broadcast`
        Queue-->>WS: 在线玩家接收事件
        WS-->>Client: 广播玩家离开信息
    end
```
