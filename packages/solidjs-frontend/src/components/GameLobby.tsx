import { createSignal, onCleanup, createEffect } from 'solid-js'
import { WebSocketService } from '../services/WebSocketService'
import type { RoomInfo } from '../services/WebSocketService'

interface GameLobbyProps {
  token: string
  onJoinGame: (gameId: string) => void
}

export function GameLobby(props: GameLobbyProps) {
  const [wsService, setWsService] = createSignal<WebSocketService | null>(null)
  const [roomInfo, setRoomInfo] = createSignal<RoomInfo | null>(null)
  const [isConnected, setIsConnected] = createSignal(false)
  const [gameId, setGameId] = createSignal('')
  const [playerId] = createSignal(Math.random().toString(36).substr(2, 9))
  const [playerName, setPlayerName] = createSignal('Player')
  const [error, setError] = createSignal('')

  createEffect(() => {
    const ws = new WebSocketService(props.token)
    setWsService(ws)

    // Setup message handlers
    ws.onMessage('roomInfo', (message) => {
      const roomInfo = message as RoomInfo
      console.log('Room info received:', roomInfo)
      setRoomInfo(roomInfo)
    })

    ws.onMessage('waiting', (message) => {
      console.log('Waiting for players:', message)
    })

    ws.onMessage('gameStart', (message) => {
      console.log('Game starting:', message)
      props.onJoinGame(roomInfo()?.roomId || gameId())
    })

    ws.onMessage('connection', (message) => {
      console.log('Connection status:', message)
      if (message.status === 'connected') {
        setIsConnected(true)
      } else {
        setError(message.message || 'Connection failed')
      }
    })

    ws.onMessage('error', (message) => {
      setError(message.error || 'Unknown error occurred')
    })

    ws.onConnection({
      onOpen: () => {
        setIsConnected(true)
        setError('')
      },
      onClose: () => {
        setIsConnected(false)
      },
      onError: () => {
        setError('WebSocket connection failed')
        setIsConnected(false)
      }
    })

    // Connect to WebSocket
    ws.connect().catch(() => {
      setError('Failed to connect to game server')
    })

    onCleanup(() => {
      ws.disconnect()
    })
  })

  const handleCreateGame = () => {
    const ws = wsService()
    if (!ws) return

    const newGameId = `game-${Date.now()}`
    setGameId(newGameId)
    ws.createGame(newGameId)
  }

  const handleJoinGame = () => {
    const ws = wsService()
    if (!ws || !gameId()) return

    ws.joinRoom(gameId(), playerId(), playerName())
  }

  const handleForceStart = () => {
    const ws = wsService()
    const room = roomInfo()
    if (!ws || !room) return

    ws.forceStart(room.roomId, playerId())
  }

  return (
    <div class="lobby-container">
      <h1>游戏大厅</h1>
      
      <div class="connection-status">
        状态: {isConnected() ? '已连接' : '未连接'}
      </div>

      {error() && (
        <div class="error-message">{error()}</div>
      )}

      <div class="player-info">
        <label>
          玩家名称:
          <input
            type="text"
            value={playerName()}
            onInput={(e) => setPlayerName(e.currentTarget.value)}
          />
        </label>
      </div>

      <div class="game-actions">
        <button onClick={handleCreateGame} disabled={!isConnected()}>
          创建游戏
        </button>

        <div class="join-game">
          <input
            type="text"
            placeholder="游戏ID"
            value={gameId()}
            onInput={(e) => setGameId(e.currentTarget.value)}
          />
          <button 
            onClick={handleJoinGame} 
            disabled={!isConnected() || !gameId()}
          >
            加入游戏
          </button>
        </div>
      </div>

      {roomInfo() && (
        <div class="room-info">
          <h3>房间信息</h3>
          <p>房间ID: {roomInfo()!.roomId}</p>
          <p>游戏模式: {roomInfo()!.gameMode.name}</p>
          <p>玩家数量: {roomInfo()!.players.length} / {roomInfo()!.gameMode.maxPlayers}</p>
          
          <div class="players-list">
            <h4>玩家列表:</h4>
            {roomInfo()!.players.map((player) => (
              <div class="player-item">
                <span>{player.username}</span>
                <span class={`status ${player.status}`}>
                  {player.status === 'online' ? '在线' : '离线'}
                </span>
                {player.isReady && <span class="ready">准备就绪</span>}
              </div>
            ))}
          </div>

          {roomInfo()!.players.length >= roomInfo()!.gameMode.minPlayers && (
            <button onClick={handleForceStart} class="force-start-btn">
              强制开始游戏
            </button>
          )}
        </div>
      )}
    </div>
  )
}