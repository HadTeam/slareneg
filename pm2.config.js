module.exports = {
  apps: [
    {
      name: 'go-backend',
      cwd: './packages/server',
      script: 'go',
      args: 'run .',
      interpreter: 'none',
      env: {
        PORT: '8080'
      }
    },
    {
      name: 'vite-frontend',
      cwd: './packages/ui-playground',
      script: 'pnpm',
      args: 'run dev',
      env: {
        PORT: '5173'
      }
    },
    {
      name: 'caddy',
      script: 'caddy',
      args: 'run',
      interpreter: 'none'
    }
  ]
};
