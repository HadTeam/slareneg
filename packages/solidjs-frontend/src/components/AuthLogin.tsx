import { createSignal } from 'solid-js'

interface AuthLoginProps {
  onLogin: (token: string) => void
}

export function AuthLogin(props: AuthLoginProps) {
  const [username, setUsername] = createSignal('')
  const [password, setPassword] = createSignal('')
  const [isLoading, setIsLoading] = createSignal(false)
  const [error, setError] = createSignal('')

  const handleLogin = async (e: Event) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username(),
          password: password(),
        }),
      })

      const data = await response.json()

      if (data.status === 'success') {
        props.onLogin(data.token)
      } else {
        setError(data.message || '登录失败')
      }
    } catch (err) {
      setError('网络错误，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  const handleRegister = async (e: Event) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username(),
          password: password(),
        }),
      })

      const data = await response.json()

      if (data.status === 'success') {
        props.onLogin(data.token)
      } else {
        setError(data.message || '注册失败')
      }
    } catch (err) {
      setError('网络错误，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div class="auth-container">
      <h1>Slareneg 游戏</h1>
      <form class="auth-form">
        <div class="form-group">
          <label for="username">用户名:</label>
          <input
            id="username"
            type="text"
            value={username()}
            onInput={(e) => setUsername(e.currentTarget.value)}
            disabled={isLoading()}
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码:</label>
          <input
            id="password"
            type="password"
            value={password()}
            onInput={(e) => setPassword(e.currentTarget.value)}
            disabled={isLoading()}
            required
          />
        </div>

        {error() && (
          <div class="error-message">{error()}</div>
        )}

        <div class="button-group">
          <button 
            type="submit" 
            onClick={handleLogin}
            disabled={isLoading() || !username() || !password()}
          >
            {isLoading() ? '登录中...' : '登录'}
          </button>
          
          <button 
            type="button" 
            onClick={handleRegister}
            disabled={isLoading() || !username() || !password()}
          >
            {isLoading() ? '注册中...' : '注册'}
          </button>
        </div>
      </form>
    </div>
  )
}