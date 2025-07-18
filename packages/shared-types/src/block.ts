// Type aliases for numeric types
export type Num = number; // uint16 in Go
export type Owner = number; // uint16 in Go
export type Name = string;

// Meta interface
export interface Meta {
  name: Name;
  description: string;
}

// AllowMove interface
export interface AllowMove {
  from: boolean;
  to: boolean;
  reason: string; // debug info
}

// Block interface
export interface Block {
  num(): Num;
  owner(): Owner;
  
  // Round Events
  roundStart(roundNum: number): void;
  roundEnd(roundNum: number): void;
  
  // Move related
  allowMove(): AllowMove;
  moveFrom(num: Num): Num;
  // MoveTo returns a new block to replace this place
  moveTo(num: Num, owner: Owner): Block;
  
  fog(isOwner: boolean, isSight: boolean): Block;
  
  meta(): Meta;
}
