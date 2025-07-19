import type { Block, Meta, AllowMove, Num, Owner } from '@slareneg/shared-types';

// Mock implementation of Block interface
export class MockBlock implements Block {
  private _num: Num;
  private _owner: Owner;
  private _meta: Meta;

  constructor(num: Num, owner: Owner, meta: Meta) {
    this._num = num;
    this._owner = owner;
    this._meta = meta;
  }

  num(): Num {
    return this._num;
  }

  owner(): Owner {
    return this._owner;
  }

  roundStart(_roundNum: number): void {
    // Mock implementation
  }

  roundEnd(_roundNum: number): void {
    // Mock implementation
  }

  allowMove(): AllowMove {
    return {
      from: true,
      to: true,
      reason: "mock allow move"
    };
  }

  moveFrom(_num: Num): Num {
    return this._num;
  }

  moveTo(num: Num, owner: Owner): Block {
    return new MockBlock(num, owner, this._meta);
  }

  fog(_isOwner: boolean, _isSight: boolean): Block {
    return this;
  }

  meta(): Meta {
    return this._meta;
  }
}

// Helper function to create blocks by name
export function createMockBlock(name: string, num: Num, owner: Owner): Block {
  const meta: Meta = {
    name,
    description: `Mock ${name} block`
  };
  return new MockBlock(num, owner, meta);
}

// Helper function to create a test board
export function createTestBoard(width: number, height: number): Block[][] {
  const blocks: Block[][] = [];
  
  for (let y = 0; y < height; y++) {
    const row: Block[] = [];
    for (let x = 0; x < width; x++) {
      const index = y * width + x;
      
      // Create different block types based on position
      let blockName = 'blank';
      if (x === 0 && y === 0) blockName = 'castle';
      else if (x === width - 1 && y === height - 1) blockName = 'castle';
      else if (x === Math.floor(width / 2) && y === Math.floor(height / 2)) blockName = 'king';
      else if ((x + y) % 5 === 0) blockName = 'mountain';
      else if ((x + y) % 3 === 0) blockName = 'soldier';
      
      const owner = blockName === 'blank' || blockName === 'mountain' ? 0 : (x < width / 2 ? 1 : 2);
      row.push(createMockBlock(blockName, index, owner));
    }
    blocks.push(row);
  }
  
  return blocks;
}
