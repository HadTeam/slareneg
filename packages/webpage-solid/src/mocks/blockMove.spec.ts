import { describe, it, expect } from 'vitest';
import { createMockBlock } from './mockBlocks';
import type { Block, Owner } from '@slareneg/shared-types';

describe('MockBlock', () => {
  describe('allowMove', () => {
    it('should always allow moves for mock blocks', () => {
      const block = createMockBlock('soldier', 10, 1);
      const allowMove = block.allowMove();
      
      expect(allowMove.from).toBe(true);
      expect(allowMove.to).toBe(true);
      expect(allowMove.reason).toBe('mock allow move');
    });
  });

  describe('moveFrom', () => {
    it('should return the current block number', () => {
      const block = createMockBlock('soldier', 42, 1);
      const result = block.moveFrom(50);
      
      expect(result).toBe(42);
    });
  });

  describe('moveTo', () => {
    it('should create a new block with the specified number and owner', () => {
      const originalBlock = createMockBlock('castle', 10, 1);
      const newOwner: Owner = 2;
      const newNum = 25;
      
      const movedBlock = originalBlock.moveTo(newNum, newOwner);
      
      expect(movedBlock.num()).toBe(newNum);
      expect(movedBlock.owner()).toBe(newOwner);
      expect(movedBlock.meta().name).toBe('castle'); // Should preserve the block type
    });

    it('should preserve metadata when moving', () => {
      const originalBlock = createMockBlock('mountain', 5, 0);
      const movedBlock = originalBlock.moveTo(10, 1);
      
      expect(movedBlock.meta().name).toBe('mountain');
      expect(movedBlock.meta().description).toBe('Mock mountain block');
    });
  });

  describe('owner', () => {
    it('should return correct owner for different block types', () => {
      const soldier = createMockBlock('soldier', 0, 1);
      const castle = createMockBlock('castle', 1, 2);
      const mountain = createMockBlock('mountain', 2, 0);
      
      expect(soldier.owner()).toBe(1);
      expect(castle.owner()).toBe(2);
      expect(mountain.owner()).toBe(0);
    });
  });

  describe('num', () => {
    it('should return the correct block number', () => {
      const block1 = createMockBlock('blank', 0, 0);
      const block2 = createMockBlock('king', 99, 1);
      const block3 = createMockBlock('soldier', 255, 2);
      
      expect(block1.num()).toBe(0);
      expect(block2.num()).toBe(99);
      expect(block3.num()).toBe(255);
    });
  });

  describe('meta', () => {
    it('should return correct metadata for each block type', () => {
      const blockTypes = ['blank', 'castle', 'king', 'mountain', 'soldier'];
      
      blockTypes.forEach(type => {
        const block = createMockBlock(type, 0, 0);
        const meta = block.meta();
        
        expect(meta.name).toBe(type);
        expect(meta.description).toBe(`Mock ${type} block`);
      });
    });
  });

  describe('fog', () => {
    it('should return the same block (no fog implementation for mocks)', () => {
      const block = createMockBlock('castle', 10, 1);
      
      const foggedBlock1 = block.fog(true, true);
      const foggedBlock2 = block.fog(false, false);
      
      expect(foggedBlock1).toBe(block);
      expect(foggedBlock2).toBe(block);
    });
  });

  describe('roundStart and roundEnd', () => {
    it('should not throw when called', () => {
      const block = createMockBlock('soldier', 0, 1);
      
      expect(() => block.roundStart(1)).not.toThrow();
      expect(() => block.roundEnd(1)).not.toThrow();
    });
  });
});

// Test helper functions for movement validation
describe('Movement validation helpers', () => {
  // Example of pure functions that could validate movement rules
  const canMoveFrom = (block: Block): boolean => {
    return block.allowMove().from;
  };

  const canMoveTo = (block: Block): boolean => {
    return block.allowMove().to;
  };

  const isValidMove = (fromBlock: Block, toBlock: Block): boolean => {
    return canMoveFrom(fromBlock) && canMoveTo(toBlock);
  };

  it('should validate movement between blocks', () => {
    const soldier = createMockBlock('soldier', 0, 1);
    const blank = createMockBlock('blank', 1, 0);
    
    expect(canMoveFrom(soldier)).toBe(true);
    expect(canMoveTo(blank)).toBe(true);
    expect(isValidMove(soldier, blank)).toBe(true);
  });

  it('should handle mountain blocks (typically immovable in real game)', () => {
    const mountain = createMockBlock('mountain', 5, 0);
    
    // In mock implementation, all moves are allowed
    expect(canMoveFrom(mountain)).toBe(true);
    expect(canMoveTo(mountain)).toBe(true);
  });
});
