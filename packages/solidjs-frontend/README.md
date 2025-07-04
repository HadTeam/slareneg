# Slareneg SolidJS Frontend

A SolidJS-based frontend for the Slareneg game with WebSocket support.

## Features

- **Authentication**: Login and registration system
- **Game Lobby**: Create and join game rooms
- **Real-time Game Board**: Interactive game board with WebSocket communication
- **Game Controls**: Move units, surrender, force start games

## Architecture

- **SolidJS**: Reactive UI framework for efficient updates
- **TypeScript**: Type-safe development
- **WebSocket**: Real-time communication with the game server
- **Vite**: Fast development and build tool

## Components

### AuthLogin
Handles user authentication with login and registration functionality.

### GameLobby  
Manages game room creation, joining, and player status.

### GameBoard
Interactive game board with:
- Real-time map display
- Unit movement controls
- Game state updates
- Turn management

### WebSocketService
Service layer for WebSocket communication:
- Connection management
- Message handling
- Game command interface

## API Integration

Connects to the Slareneg server's WebSocket endpoint at `/api/game/ws` with JWT authentication.

Supports all game commands:
- `join` - Join a game room
- `createGame` - Create a new game
- `move` - Move units on the board
- `forceStart` - Force start a game
- `surrender` - Surrender the game

## Usage

```bash
$ npm install # or pnpm install or yarn install
```

### Learn more on the [Solid Website](https://solidjs.com) and come chat with us on our [Discord](https://discord.com/invite/solidjs)

## Available Scripts

In the project directory, you can run:

### `npm run dev`

Runs the app in the development mode.<br>
Open [http://localhost:5173](http://localhost:5173) to view it in the browser.

### `npm run build`

Builds the app for production to the `dist` folder.<br>
It correctly bundles Solid in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

## Game Usage

1. **Login**: Enter username and password to authenticate
2. **Lobby**: Create a new game or join an existing one
3. **Game**: Click units to select, click adjacent cells to move
4. **Controls**: Use surrender button or leave game as needed

## Styling

The interface features:
- Modern gradient background
- Responsive grid-based game board
- Color-coded player territories
- Interactive hover effects
- Mobile-friendly design

## WebSocket Protocol

Follows the Slareneg WebSocket API specification for real-time game communication with the server.

## Deployment

Learn more about deploying your application with the [documentations](https://vite.dev/guide/static-deploy.html)
