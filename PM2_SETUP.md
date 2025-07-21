# PM2 Setup for Slareneg

This project uses PM2 to manage all services in development.

## Quick Start

```bash
# Install dependencies
pnpm install

# Start all services
pnpm start

# View logs
pnpm logs

# Check status
pnpm status

# Stop all services
pnpm stop

# Restart all services
pnpm restart

# Kill PM2 daemon
pnpm kill
```

## Services

1. **go-backend** - Go server running on port 8080
2. **vite-frontend** - Vite dev server running on port 5173
3. **caddy** - Reverse proxy running on port 800

## Access Points

- Main app: http://localhost:800
- Frontend directly: http://localhost:5173
- Backend API: http://localhost:8080/api/

## PM2 Configuration

See `pm2.config.js` for service configuration.

## Troubleshooting

If ports are already in use:
```bash
# Check what's using the ports
lsof -i :800,5173,8080

# Kill PM2 and restart
pnpm kill
pnpm start
```
