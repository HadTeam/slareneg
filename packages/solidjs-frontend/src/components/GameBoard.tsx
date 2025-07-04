import { createSignal, onCleanup, createEffect, For } from 'solid-js'
import { WebSocketService } from '../services/WebSocketService'
import type { GameStart, BlockInfo } from '../services/WebSocketService'

interface GameBoardProps {
  token: string
  gameId: string
  onLeaveGame: () => void
}

export function GameBoard(props: GameBoardProps) {
  const [wsService, setWsService] = createSignal<WebSocketService | null>(null)
  const [gameState, setGameState] = createSignal<GameStart | null>(null)
  const [currentPlayer, setCurrentPlayer] = createSignal<number>(0)
  const [turnTimeLeft, setTurnTimeLeft] = createSignal<number>(0)
  const [selectedCell, setSelectedCell] = createSignal<{x: number, y: number} | null>(null)
  const [playerId] = createSignal(Math.random().toString(36).substr(2, 9))
  const [isConnected, setIsConnected] = createSignal(false)
  const [error, setError] = createSignal('')

  createEffect(() => {
    const ws = new WebSocketService(props.token)
    setWsService(ws)

    // Setup message handlers
    ws.onMessage('gameStart', (message) => {
      const gameStart = message as GameStart
      console.log('Game started:', gameStart)
      setGameState(gameStart)
      setCurrentPlayer(gameStart.currentPlayer)
      setTurnTimeLeft(gameStart.turnTimeLeft)
    })

    ws.onMessage('newTurn', (message) => {
      console.log('New turn:', message)
      if (gameState()) {
        setGameState({
          ...gameState()!,
          map: message.map,
          turnNumber: message.turnNumber
        })
      }
      setCurrentPlayer(message.currentPlayer)
      setTurnTimeLeft(message.turnTimeLeft)
    })

    ws.onMessage('gameStateUpdate', (message) => {
      console.log('Game state update:', message)
      if (gameState()) {
        setGameState({
          ...gameState()!,
          map: message.map
        })
      }
    })

    ws.onMessage('gameEnd', (message) => {
      console.log('Game ended:', message)
      alert(`游戏结束! 获胜者: ${message.winner}`)
    })

    ws.onMessage('error', (message) => {
      setError(message.error || 'Unknown error occurred')
    })

    ws.onConnection({
      onOpen: () => {
        setIsConnected(true)
        // Join the game room
        ws.joinRoom(props.gameId, playerId(), 'Player')
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

  const getCellClass = (block: BlockInfo, x: number, y: number) => {
    const selected = selectedCell()
    const isSelected = selected && selected.x === x && selected.y === y
    
    let className = `cell ${block.type}`
    if (block.owner !== undefined) {
      className += ` player-${block.owner}`
    }
    if (isSelected) {
      className += ' selected'
    }
    return className
  }

  const getCellDisplay = (block: BlockInfo) => {
    if (block.type === 'troop' || block.type === 'city' || block.type === 'general') {
      return block.troops?.toString() || ''
    }
    if (block.type === 'mountain') {
      return '⛰️'
    }
    return ''
  }

  const handleCellClick = (x: number, y: number, block: BlockInfo) => {
    const selected = selectedCell()
    
    if (!selected) {
      // First click - select source cell
      if (block.owner === currentPlayer() && (block.type === 'troop' || block.type === 'city' || block.type === 'general')) {
        setSelectedCell({ x, y })
      }
    } else {
      // Second click - move to target cell
      if (selected.x === x && selected.y === y) {
        // Clicked same cell, deselect
        setSelectedCell(null)
      } else {
        // Attempt to move
        handleMove(selected, { x, y })
        setSelectedCell(null)
      }
    }
  }

  const handleMove = (from: {x: number, y: number}, to: {x: number, y: number}) => {
    const ws = wsService()
    if (!ws) return

    // Determine direction
    let direction: string
    if (to.y < from.y) direction = 'up'
    else if (to.y > from.y) direction = 'down'
    else if (to.x < from.x) direction = 'left'
    else if (to.x > from.x) direction = 'right'
    else return // Same position

    // Check if it's adjacent
    const dx = Math.abs(to.x - from.x)
    const dy = Math.abs(to.y - from.y)
    if (dx + dy !== 1) return // Not adjacent

    const game = gameState()
    if (!game) return

    const sourceBlock = game.map[from.y]?.[from.x]
    if (!sourceBlock || !sourceBlock.troops || sourceBlock.troops <= 1) return

    // Move half the troops (or all if only 2)
    const troops = Math.floor(sourceBlock.troops / 2)
    
    ws.move(props.gameId, playerId(), from, direction, troops)
  }

  const handleSurrender = () => {
    const ws = wsService()
    if (!ws) return

    if (confirm('确定要投降吗？')) {
      ws.surrender(props.gameId, playerId())
      props.onLeaveGame()
    }
  }

  const game = gameState()

  return (
    <div class="game-board-container">
      <div class="game-header">
        <h2>Slareneg 游戏</h2>
        <div class="game-info">
          {game && (
            <>
              <span>回合: {game.turnNumber}</span>
              <span>当前玩家: {currentPlayer()}</span>
              <span>剩余时间: {turnTimeLeft()}s</span>
            </>
          )}
        </div>
        <div class="game-controls">
          <button onClick={props.onLeaveGame}>离开游戏</button>
          <button onClick={handleSurrender} class="surrender-btn">投降</button>
        </div>
      </div>

      <div class="connection-status">
        状态: {isConnected() ? '已连接' : '未连接'}
      </div>

      {error() && (
        <div class="error-message">{error()}</div>
      )}

      {game ? (
        <div class="game-board" style={`--cols: ${game.mapWidth}; --rows: ${game.mapHeight}`}>
          <For each={game.map}>
            {(row, y) => (
              <For each={row}>
                {(block, x) => (
                  <div
                    class={getCellClass(block, x(), y())}
                    onClick={() => handleCellClick(x(), y(), block)}
                  >
                    {getCellDisplay(block)}
                  </div>
                )}
              </For>
            )}
          </For>
        </div>
      ) : (
        <div class="loading">等待游戏开始...</div>
      )}

      <div class="game-instructions">
        <h3>游戏说明:</h3>
        <ul>
          <li>点击你的单位选择，再点击相邻格子移动</li>
          <li>目标是占领所有敌人的将军</li>
          <li>数字表示该格子的兵力</li>
        </ul>
      </div>
    </div>
  )
}