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

##### 大致数据流
```mermaid
sequenceDiagram
    participant Client as 客户端 Client
    participant HTTP as HTTP Service
    participant WS as WebSocket Service
    participant Game as 游戏核心 GameCore
    participant Queue as 消息队列 Queue
    participant Cache as 缓存层 Cache<br/>(RAM/Redis)
    participant DB as 数据库 DB<br/>(Postrges)

    %% 用户认证流程
    Client->>HTTP: 注册/登录请求
    HTTP->>DB: 验证用户凭据
    DB-->>HTTP: 返回用户信息
    HTTP->>Cache: 缓存用户会话
    HTTP-->>Client: 返回JWT令牌

    %% WebSocket连接建立
    Client->>WS: 建立WebSocket连接
    WS->>Cache: 验证JWT令牌
    Cache-->>WS: 返回用户信息
    WS->>Game: 玩家加入游戏房间
    Game->>Cache: 更新房间状态
    WS-->>Client: 连接成功

    %% 游戏操作流程
    Client->>WS: 发送游戏指令(移动/攻击/投降)
    WS->>Queue: 指令入队列
    
    Queue->>Game: 消费指令
    Note over Game: 1. 验证指令合法性<br/>2. 检查玩家权限<br/>3. 判断回合状态
    
    alt 指令有效
        Game->>Cache: 更新游戏状态
        Game->>Queue: 广播状态变更
        Queue->>WS: 推送给所有玩家
        WS-->>Client: 实时状态更新
    else 指令无效
        Game-->>WS: 返回错误信息
        WS-->>Client: 显示错误提示
    end

    %% 定时任务处理
    rect rgb(240, 248, 255)
        Note over Game: 定时任务模块 Timer Module
        loop 每回合间隔
            Game->>Game: 处理资源增长
            Game->>Cache: 批量更新状态
            Game->>Queue: 广播回合结束
            Queue->>WS: 推送新回合开始
            WS-->>Client: 回合状态更新
        end
    end

    %% 游戏结束流程
    Game->>Game: 检查胜负条件
    alt 游戏结束
        Game->>DB: 保存游戏结果
        Game->>Cache: 清理游戏状态
        Game->>Queue: 广播游戏结束
        Queue->>WS: 推送结算信息
        WS-->>Client: 显示结算界面
    end

    %% 异常处理
    alt 玩家断线
        WS->>Game: 玩家离线通知
        Game->>Cache: 标记玩家状态
        Game->>Queue: 广播玩家离线
    end

    %% 数据持久化
    rect rgb(255, 248, 240)
        Note over Cache,DB: 数据同步
        Cache->>DB: 定期同步缓存数据
        DB-->>Cache: 确认同步完成
    end
```
