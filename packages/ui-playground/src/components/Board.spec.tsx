import { describe, it, expect } from 'vitest';
import { render } from '@solidjs/testing-library';
import Board from './Board';
import { createTestBoard } from '../mocks/mockBlocks';
import type { Size } from '@slareneg/shared-types';

describe('Board Component', () => {
  it('should render the board with correct dimensions', () => {
    const blocks = createTestBoard(10, 10);
    const size: Size = { width: 10, height: 10 };
    
    const { container } = render(() => (
      <Board blocks={blocks} size={size} />
    ));
    
    // Check that the grid container exists
    const gridContainer = container.querySelector('[style*="grid-template-columns"]');
    expect(gridContainer).toBeTruthy();
    
    // Check grid dimensions in style
    const style = gridContainer?.getAttribute('style');
    expect(style).toContain('grid-template-columns: repeat(10, 50px)');
    expect(style).toContain('grid-template-rows: repeat(10, 50px)');
  });

  it('should render correct number of cells', () => {
    const blocks = createTestBoard(5, 5);
    const size: Size = { width: 5, height: 5 };
    
    render(() => <Board blocks={blocks} size={size} />);
    
    // Count cells by their 50x50 dimensions
    const cells = document.querySelectorAll('div[style*="width: 50px"][style*="height: 50px"]');
    expect(cells.length).toBe(25); // 5x5 grid
  });

  it('should render different block types with Unicode characters', () => {
    const blocks = createTestBoard(10, 10);
    const size: Size = { width: 10, height: 10 };
    
    const { container } = render(() => <Board blocks={blocks} size={size} />);
    
    // Check that we have some of each block type
    const content = container.textContent || '';
    
    // Mountain (▲)
    expect(content).toContain('▲');
    
    // Castle (🏛)
    expect(content).toContain('🏛');
    
    // King (♔)
    expect(content).toContain('♔');
    
    // Soldier (◆)
    expect(content).toContain('◆');
  });

  it('should apply owner colors to cells', () => {
    const blocks = createTestBoard(4, 4);
    const size: Size = { width: 4, height: 4 };
    
    const { container } = render(() => <Board blocks={blocks} size={size} />);
    
    // Get all cells
    const cells = container.querySelectorAll('div[style*="width: 50px"]');
    
    // Check that we have cells with different background colors
    const styles = Array.from(cells).map(cell => cell.getAttribute('style') || '');
    
    // Should have some cells with player colors or block colors
    const hasColoredCells = styles.some(style => 
      style.includes('background:') && !style.includes('background: transparent')
    );
    
    expect(hasColoredCells).toBe(true);
  });

  it('should handle zoom functionality', () => {
    const blocks = createTestBoard(3, 3);
    const size: Size = { width: 3, height: 3 };
    
    const { container } = render(() => <Board blocks={blocks} size={size} />);
    
    // Check that transform div exists for zoom/pan
    const transformDiv = container.querySelector('[style*="transform"]');
    expect(transformDiv).toBeTruthy();
    
    // Should have initial scale
    const style = transformDiv?.getAttribute('style') || '';
    expect(style).toContain('scale(1)');
  });

  it('should support cell selection', () => {
    const blocks = createTestBoard(3, 3);
    const size: Size = { width: 3, height: 3 };
    
    const { container } = render(() => <Board blocks={blocks} size={size} />);
    
    // Get cells that can be clicked
    const cells = container.querySelectorAll('div[style*="width: 50px"][style*="cursor: pointer"]');
    expect(cells.length).toBeGreaterThan(0);
    
    // Cells should have click handlers (indicated by cursor: pointer)
    const firstCell = cells[0] as HTMLElement;
    const style = firstCell.getAttribute('style') || '';
    expect(style).toContain('cursor: pointer');
  });
});
