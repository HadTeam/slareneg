import { describe, it, expect } from 'vitest';
import { generateMockMap } from './mockMap';
import type { Block } from '@slareneg/shared-types';

describe('generateMockMap', () => {
  it('should generate a 20x20 map', () => {
    const map = generateMockMap();
    
    // Check dimensions
    expect(map).toHaveLength(20);
    map.forEach(row => {
      expect(row).toHaveLength(20);
    });
  });

  it('should contain valid block types', () => {
    const map = generateMockMap();
    const validTypes = ['blank', 'castle', 'king', 'mountain', 'soldier'];
    
    map.forEach(row => {
      row.forEach(block => {
        expect(validTypes).toContain(block.meta().name);
      });
    });
  });

  it('should have reasonable distribution of block types', () => {
    // Generate multiple maps to get statistical averages
    const sampleSize = 10;
    const counts = {
      blank: 0,
      castle: 0,
      king: 0,
      mountain: 0,
      soldier: 0
    };
    
    for (let i = 0; i < sampleSize; i++) {
      const map = generateMockMap();
      map.forEach(row => {
        row.forEach(block => {
          const blockType = block.meta().name;
          counts[blockType as keyof typeof counts]++;
        });
      });
    }
    
    const totalBlocks = 20 * 20 * sampleSize;
    const expectedAverage = totalBlocks / 5; // 5 block types
    
    // Each block type should appear roughly 20% of the time (±10%)
    Object.entries(counts).forEach(([blockType, count]) => {
      const percentage = (count / totalBlocks) * 100;
      expect(percentage).toBeGreaterThan(10); // At least 10%
      expect(percentage).toBeLessThan(30); // At most 30%
    });
  });

  it('should assign correct owners to blocks', () => {
    const map = generateMockMap();
    
    map.forEach(row => {
      row.forEach(block => {
        const blockType = block.meta().name;
        const owner = block.owner();
        
        if (blockType === 'blank' || blockType === 'mountain') {
          // Blank and mountain blocks should have no owner (0)
          expect(owner).toBe(0);
        } else {
          // Other blocks should be owned by player 1 or 2
          expect([1, 2]).toContain(owner);
        }
      });
    });
  });

  it('should generate blocks with correct sequential numbers', () => {
    const map = generateMockMap();
    let expectedNum = 0;
    
    map.forEach((row, y) => {
      row.forEach((block, x) => {
        expect(block.num()).toBe(expectedNum);
        expectedNum++;
      });
    });
  });

  it('should have at least some mountains and castles in each generated map', () => {
    const map = generateMockMap();
    let mountainCount = 0;
    let castleCount = 0;
    
    map.forEach(row => {
      row.forEach(block => {
        const blockType = block.meta().name;
        if (blockType === 'mountain') mountainCount++;
        if (blockType === 'castle') castleCount++;
      });
    });
    
    // Expect at least 5% of blocks to be mountains (20 blocks out of 400)
    expect(mountainCount).toBeGreaterThan(20);
    // Expect at least 5% of blocks to be castles
    expect(castleCount).toBeGreaterThan(20);
  });
});
