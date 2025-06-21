## WebSocket

### 建立连接
**地址**: `ws://host/api/game/ws`
**连接参数**: 
```
Authorization: Bearer <JWT_TOKEN>
```

### 连接状态消息
服务器在连接建立后发送：
```json
{
  "type": "connection",
  "status": "connected|failed",
  "playerId": "<uint>",
  "message": "<string>"
}
```

## 客户端指令 (Client -> Server)

### 移动指令
```
Move <x> <y> <direction> <amount>
```
- **x**: 起始列坐标 (uint)
- **y**: 起始行坐标 (uint)  
- **direction**: 移动方向
  - `up`: 目标Y坐标 = y-1 (向上)
  - `down`: 目标Y坐标 = y+1 (向下)  
  - `left`: 目标X坐标 = x-1 (向左)
  - `right`: 目标X坐标 = x+1 (向右)
- **amount**: 移动数量类型
  - `0`: 全部移动
  - `1`: 移动一半

### 强制开始
```
ForceStart <enable>
```
- **enable**: `true`激活 | `false`取消

### 投降
```
Surrender
```

### 加入房间
```
JoinRoom <roomId>
```
- **roomId**: 房间ID，空表示匹配任意房间

## 服务器消息 (Server -> Client)

### 房间信息
```json
{
  "type": "roomInfo",
  "roomId": "<string>",
  "players": [
    {
      "id": "<uint>",
      "username": "<string>",
      "teamId": "<uint>",
      "status": "online|offline",
      "forceStart": "<bool>",
      "isReady": "<bool>"
    }
  ],
  "gameMode": {
    "name": "<string>",
    "maxPlayers": "<uint>",
    "minPlayers": "<uint>",
    "turnDuration": "<uint>" // 秒
  }
}
```

### 等待状态
```json
{
  "type": "waiting",
  "message": "等待其他玩家加入...",
  "currentPlayers": "<uint>",
  "requiredPlayers": "<uint>"
}
```

### 游戏开始
```json
{
  "type": "gameStart",
  "mapWidth": "<uint>",
  "mapHeight": "<uint>",
  "map": "[][]<BlockInfo>",
  "turnNumber": 1,
  "currentPlayer": "<uint>",
  "turnTimeLeft": "<uint>"
}
```

### 新回合
```json
{
  "type": "newTurn",
  "turnNumber": "<uint>",
  "map": "[][]<BlockInfo>",
  "currentPlayer": "<uint>",
  "turnTimeLeft": "<uint>",
  "lastActions": [
    {
      "playerId": "<uint>",
      "action": "<string>",
      "result": "success|failed|conflict"
    }
  ]
}
```

### 游戏状态更新
```json
{
  "type": "gameUpdate",
  "map": "[][]<BlockInfo>",
  "changedBlocks": [
    {
      "x": "<uint>",
      "y": "<uint>",
      "blockInfo": "<BlockInfo>"
    }
  ]
}
```

### 游戏结束
```json
{
  "type": "gameEnd",
  "winner": {
    "teamId": "<uint>",
    "players": ["<string>"]
  },
  "gameStats": {
    "duration": "<uint>", // 秒
    "totalTurns": "<uint>",
    "playerStats": [
      {
        "playerId": "<uint>",
        "territoryCaptured": "<uint>",
        "unitsLost": "<uint>",
        "actionsPerformed": "<uint>"
      }
    ]
  }
}
```

### 错误消息
```json
{
  "type": "error",
  "code": "<string>",
  "message": "<string>",
  "details": "<object>" // 可选
}
```

### 玩家状态变更
```json
{
  "type": "playerUpdate",
  "playerId": "<uint>",
  "status": "joined|left|reconnected|timeout",
  "message": "<string>"
}
```

## 数据结构定义

### BlockInfo
```typescript
type BlockInfo = [
  blockType: number,  // 0:空地 1:障碍 2:城市 3:山脉等
  ownerId: number,    // 0:中立 其他:玩家ID
  unitCount: number   // 单位数量
]
```

### 错误码定义
- `INVALID_TOKEN`: JWT令牌无效
- `GAME_NOT_FOUND`: 游戏房间不存在  
- `INVALID_MOVE`: 无效移动
- `NOT_YOUR_TURN`: 不是该玩家回合
- `GAME_ENDED`: 游戏已结束
- `PLAYER_LIMIT_REACHED`: 房间人数已满
- `INSUFFICIENT_UNITS`: 单位数量不足
- `INVALID_POSITION`: 坐标无效
