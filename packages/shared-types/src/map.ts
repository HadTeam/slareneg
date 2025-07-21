import { Block, Owner } from './block';

// Size interface
export interface Size {
  width: number; // uint16 in Go
  height: number; // uint16 in Go
}

// Pos interface
export interface Pos {
  x: number; // uint16 in Go
  y: number; // uint16 in Go
}

// Info interface
export interface Info {
  id: string;
  name: string;
  desc: string;
}

// Type aliases for clarity
export type Sight = boolean[][]; // 2D array for visibility, true if visible
export type Blocks = Block[][]; // 2D array for blocks

// Map interface
export interface Map {
  isEmpty(): boolean;
  block(pos: Pos): Block | null; // Returns null instead of error in TS
  blocks(): Blocks;
  setBlock(pos: Pos, b: Block): void; // No error return in TS
  setBlocks(blocks: Blocks): void; // Set all blocks at once
  size(): Size;
  info(): Info;
  
  roundStart(roundNum: number): void;
  roundEnd(roundNum: number): void;
  
  fog(owner: Owner[], sight: Sight): void; // No error return in TS
}

// Helper functions
export function sizeToString(s: Size): string {
  return `${s.width}x${s.height}`;
}

export function posToString(p: Pos): string {
  return `Pos(${p.x},${p.y})`;
}

export function isPosValid(s: Size, p: Pos): boolean {
  return p.x >= 1 && p.x <= s.width && p.y >= 1 && p.y <= s.height;
}

export function infoToString(i: Info): string {
  return `Info(#${i.id}, ${i.name}, Desc: ${i.desc})`;
}
