// WebSocket message types based on the API documentation
export interface ClientMessage {
  type: string
  gameId?: string
  payload?: any
}

export interface JoinPayload {
  playerId: string
  playerName: string
}

export interface MovePayload {
  playerId: string
  from: { x: number; y: number }
  direction: 'up' | 'down' | 'left' | 'right'
  troops: number
}

export interface ServerMessage {
  type: string
  [key: string]: any
}

export interface RoomInfo {
  type: 'roomInfo'
  roomId: string
  players: Player[]
  gameMode: GameMode
}

export interface Player {
  id: number
  username: string
  teamId: number
  status: 'online' | 'offline'
  forceStart: boolean
  isReady: boolean
}

export interface GameMode {
  name: string
  maxPlayers: number
  minPlayers: number
  turnDuration: number
}

export interface GameStart {
  type: 'gameStart'
  mapWidth: number
  mapHeight: number
  map: BlockInfo[][]
  turnNumber: number
  currentPlayer: number
  turnTimeLeft: number
}

export interface BlockInfo {
  type: 'empty' | 'troop' | 'city' | 'general' | 'mountain'
  owner?: number
  troops?: number
}

export class WebSocketService {
  private ws: WebSocket | null = null
  private messageHandlers: Map<string, (message: ServerMessage) => void> = new Map()
  private connectionHandlers: {
    onOpen?: () => void
    onClose?: () => void
    onError?: (error: Event) => void
  } = {}

  constructor(_token: string) {
    // Store token for future use if needed
    console.log('WebSocket service initialized with token')
  }

  connect(): Promise<boolean> {
    return new Promise((resolve, reject) => {
      try {
        // Use ws:// for development, wss:// for production
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/api/game/ws`
        
        this.ws = new WebSocket(wsUrl)

        this.ws.onopen = () => {
          console.log('WebSocket connected')
          this.connectionHandlers.onOpen?.()
          resolve(true)
        }

        this.ws.onclose = () => {
          console.log('WebSocket disconnected')
          this.connectionHandlers.onClose?.()
        }

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error)
          this.connectionHandlers.onError?.(error)
          reject(error)
        }

        this.ws.onmessage = (event) => {
          try {
            const message: ServerMessage = JSON.parse(event.data)
            console.log('Received message:', message)
            
            const handler = this.messageHandlers.get(message.type)
            if (handler) {
              handler(message)
            }
          } catch (err) {
            console.error('Failed to parse message:', err)
          }
        }
      } catch (error) {
        reject(error)
      }
    })
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  sendMessage(message: ClientMessage) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
      console.log('Sent message:', message)
    } else {
      console.error('WebSocket not connected')
    }
  }

  // Game specific methods
  joinRoom(gameId: string, playerId: string, playerName: string) {
    this.sendMessage({
      type: 'join',
      gameId,
      payload: { playerId, playerName }
    })
  }

  createGame(gameId: string) {
    this.sendMessage({
      type: 'createGame',
      gameId,
      payload: {}
    })
  }

  move(gameId: string, playerId: string, from: { x: number; y: number }, direction: string, troops: number) {
    this.sendMessage({
      type: 'move',
      gameId,
      payload: { playerId, from, direction, troops }
    })
  }

  forceStart(gameId: string, playerId: string) {
    this.sendMessage({
      type: 'forceStart',
      gameId,
      payload: { playerId }
    })
  }

  surrender(gameId: string, playerId: string) {
    this.sendMessage({
      type: 'surrender',
      gameId,
      payload: { playerId }
    })
  }

  // Event handlers
  onMessage(type: string, handler: (message: ServerMessage) => void) {
    this.messageHandlers.set(type, handler)
  }

  onConnection(handlers: {
    onOpen?: () => void
    onClose?: () => void
    onError?: (error: Event) => void
  }) {
    this.connectionHandlers = handlers
  }
}