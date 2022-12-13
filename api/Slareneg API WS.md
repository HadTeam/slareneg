> **约定**
> 开头大写的量，允许在全局被引用

## Connection
addr:  `/api/game/ws`

## Client-Side
大概格式为 `<Instruction> <...Args>`
以空格为指令内分隔符，`\n` 为指令分隔符
### Move
`Move <x> <y> <towards>`
- x uint: 列数
- y uint: 行数
- towards string:
  - "up": y 减小的方向
  - "down": y 增大的方向
  - "left": x 减小的方向
  - "right": x 增大的方向
### ForceStart
`ForceStart <status>`
- status bool:
  - true: 激活
  - false: 取消激活
### Surrender
`Surrender`

## Server-Side
返回 Base64 加密的 JSON 文本
### Start
```json
{
  "action": "start",
  "mapWidth": <uint>,
  "mapHeight": <uint>,
  "map": [mapHeight][mapWidth]<BlockTypeId(uint)>
}
```
#### BlockInfo(uint)
`<BlockTypeInfo{2}><Number>`
##### BlockTypeInfo(uint)
```
00=>blank
01=>king
02=>mountain
03=>city
```

### Wait
```json
{
  "action": "wait"
  "players": []<PlayerBaseInfo>,
  "minNumber": <uint>
}
```
#### PlayerBaseInfo(struct)
```json
{
  "name": <string>,
  "playerId": <uint>
  "forceStart": <bool>,
  "teamId": <TeamId>
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
  "map": [mapHeight][mapWidth]<BlockInfo(uint)>,
  "mapCover": [mapHeight][mapWidth]<TeamId(uint)>,
  "scoreBoard": {
    <playerId>: {
      "number": <uint>,
      "place": <uint>
    }
  }
}
```
#### TeamId(uint)
```
0=>neutral
```

