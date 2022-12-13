## User
### 注册
addr: `/api/user/register?username={username}`
cookie: `token: <token>`

### 恢复

addr: `/api/user/recover?mail={mail_addr}`
ret: `{ "status": "<Status>" }`

#### Status
- success
- failed

### 绑定邮箱

addr: `/api/user/bind?mail={mail_addr}`
ret: `{ "status": "<Status>" }`

