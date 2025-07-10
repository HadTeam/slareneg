import { createSignal, Show } from 'solid-js'
import './App.css'
import { AuthLogin } from './components/AuthLogin'
import { GameLobby } from './components/GameLobby'
import { GameBoard } from './components/GameBoard'

type AppState = 'login' | 'lobby' | 'game'

function App() {
  const [appState, setAppState] = createSignal<AppState>('login')
  const [userToken, setUserToken] = createSignal<string>('')
  const [gameId, setGameId] = createSignal<string>('')

  const handleLogin = (token: string) => {
    setUserToken(token)
    setAppState('lobby')
  }

  const handleJoinGame = (id: string) => {
    setGameId(id)
    setAppState('game')
  }

  const handleLeaveGame = () => {
    setGameId('')
    setAppState('lobby')
  }

  return (
    <div class="app">
      <Show when={appState() === 'login'}>
        <AuthLogin onLogin={handleLogin} />
      </Show>
      <Show when={appState() === 'lobby'}>
        <GameLobby 
          token={userToken()} 
          onJoinGame={handleJoinGame} 
        />
      </Show>
      <Show when={appState() === 'game'}>
        <GameBoard 
          token={userToken()} 
          gameId={gameId()} 
          onLeaveGame={handleLeaveGame}
        />
      </Show>
    </div>
  )
}

export default App
