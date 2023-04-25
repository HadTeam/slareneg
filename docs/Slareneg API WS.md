> **约定**
> 开头大写的量，允许在全局被引用

## Connection
addr:  `/api/game/ws`

## Client-Side
大概格式为 `<Instruction> <...Args>`
以空格为指令内分隔符，`\n` 为指令分隔符
### Move
`Move <x> <y> <towards> <number>`
- x uint: 列数
- y uint: 行数
- towards string:
  - "up": y 减小的方向
  - "down": y 增大的方向
  - "left": x 减小的方向
  - "right": x 增大的方向
- number: 移动的数量
  - 0: 全部移动（对应原版 `is50: false`）
  - 1: 移动一半（对应原版 `is50: true`）
### ForceStart
`ForceStart <status>`
- status bool:
  - true: 激活
  - false: 取消激活
### Surrender
`Surrender`

## Server-Side

[//]: # (返回 Base64 加密的 JSON 文本)
### Start
```json
{
  "action": "start",
  "mapWidth": "<uint>",
  "mapHeight": "<uint>",
  "map": [][]<BlockInfo>
}
```
#### BlockInfo([]uint)
`[<Block(Type)Id>, <OwnerId>, <Number>]`

### Wait
```json
{
  "action": "wait"
}
```

### Info
```json
{
  "players": []<PlayerBaseInfo>,
  "mode": {
    "MaxUserNum": <uint>,
    "MinUserNum": <uint>,
    "NameStr": <string>
  }
}
```

#### PlayerInfo(struct)
```json
{
  "name": <string>,
  "id": <uint>,
  "forceStart": <bool>,
  "teamId": <TeamId>,
  "status": <bool>
}
```

### End
```json
{
  "action": "end",
  "winnerTeam": <TeamId>
}
```

### NewTurn
```json
{
  "action": "newTurn",
  "turnNumber": <uint>,
  "map": [][]<BlockInfo>
}
```
