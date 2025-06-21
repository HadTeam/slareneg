## 认证系统 (Authentication)

### 用户注册
**地址**: `POST /api/auth/register`
**请求体**:
```json
{
  "username": "<string>",
  "password": "<string>",
  "email": "<string>" // 可选
}
```
**响应**:
```json
{
  "status": "success|failed",
  "token": "<JWT_TOKEN>", // 成功时返回
  "message": "<string>"   // 失败时返回错误信息
}
```

### 用户登录
**地址**: `POST /api/auth/login`
**请求体**:
```json
{
  "username": "<string>",
  "password": "<string>"
}
```
**响应**:
```json
{
  "status": "success|failed",
  "token": "<JWT_TOKEN>",
  "userInfo": {
    "id": "<uint>",
    "username": "<string>",
    "email": "<string>"
  }
}
```

### 密码恢复
**地址**: `POST /api/auth/recover`
**请求体**:
```json
{
  "email": "<string>"
}
```
**响应**:
```json
{
  "status": "success|failed",
  "message": "<string>"
}
```

### 绑定邮箱
**地址**: `POST /api/auth/bind-email`
**请求头**: `Authorization: Bearer <JWT_TOKEN>`
**请求体**:
```json
{
  "email": "<string>"
}
```
**响应**:
```json
{
  "status": "success|failed",
  "message": "<string>"
}
```