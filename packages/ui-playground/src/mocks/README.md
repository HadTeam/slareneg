# Mock Map Implementation

## Overview
This directory contains mock implementations for testing the Slareneg game board.

## Files

### mockMap.ts
- Generates a random 20×20 game board matrix
- Uses block types: blank, castle, king, mountain, soldier
- Randomly assigns owners (0 for neutral, 1 or 2 for players)
- Mountains and blank cells are always neutral (owner = 0)

### mockBlocks.ts
- Contains the Block interface implementation
- Provides helper functions to create test boards
- Implements all required block methods

## UI Features Implemented

1. **Full-screen board as background**
   - Dark theme (#2a2a2a background)
   - Board centered on screen
   - Grid gaps for visual separation

2. **Improved cell visualization**
   - Outline-style Unicode icons
   - Stacked icon + number display
   - Color-coded by owner (blue/red for players, gray for neutral)
   - Hover effects with scale animation

3. **Cell selection**
   - Click to select/deselect cells
   - White border for selected cells
   - Hover highlights with owner color

4. **Pan and zoom**
   - Mouse wheel to zoom (0.3x to 3x)
   - Click and drag to pan the board
   - Smooth transitions

## Block Types
- **Blank**: Empty cell (no icon)
- **Castle**: 🏛 (Greek building)
- **King**: ♔ (Chess king)
- **Mountain**: ▲ (Triangle)
- **Soldier**: ◆ (Diamond)
