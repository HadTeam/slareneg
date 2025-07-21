# @slareneg/shared-types

Shared TypeScript type definitions for the slareneg game. These types are translations from the Go backend types to ensure type safety across the frontend and testing environments.

## Installation

From within a workspace package:

```bash
pnpm add @slareneg/shared-types
```

## Usage

```typescript
import { Block, Map, Size, Pos, Owner, Num } from '@slareneg/shared-types';

// Use the types in your code
const position: Pos = { x: 5, y: 10 };
const mapSize: Size = { width: 20, height: 20 };
```

## Available Types

### Block Types
- `Num` - Numeric type for block numbers
- `Owner` - Numeric type for block owners
- `Name` - String type for names
- `Meta` - Block metadata interface
- `AllowMove` - Movement permission interface
- `Block` - Main block interface

### Map Types
- `Size` - Map size interface (width, height)
- `Pos` - Position interface (x, y)
- `Info` - Map information interface
- `Sight` - 2D boolean array for visibility
- `Blocks` - 2D array of Block interfaces
- `Map` - Main map interface

### Helper Functions
- `sizeToString(s: Size): string` - Convert size to string format
- `posToString(p: Pos): string` - Convert position to string format
- `isPosValid(s: Size, p: Pos): boolean` - Check if position is valid within size
- `infoToString(i: Info): string` - Convert info to string format

## Notes

These types are direct translations from the Go backend types with some adjustments for TypeScript:
- Go's `uint16` types are represented as `number` in TypeScript
- Error returns are replaced with `null` or `void` where appropriate
- Method names follow TypeScript conventions (camelCase)
