import { For, createSignal, onMount, onCleanup } from 'solid-js';
import type { Blocks, Size } from '@slareneg/shared-types';
import Cell from './Cell';
import { createMockBlock } from '../mocks/mockBlocks';

interface BoardProps {
  blocks?: Blocks;
  size: Size;
  onBoardRef?: (ref: { fitToView: () => void }) => void;
}

function Board(props: BoardProps) {
  const { blocks, size, onBoardRef } = props;
  const [scale, setScale] = createSignal(1);
  const [translateX, setTranslateX] = createSignal(0);
  const [translateY, setTranslateY] = createSignal(0);
  const [isPanning, setIsPanning] = createSignal(false);
  const [startX, setStartX] = createSignal(0);
  const [startY, setStartY] = createSignal(0);
  const [selectedCell, setSelectedCell] = createSignal<{x: number, y: number} | null>(null);

  let containerRef: HTMLDivElement | undefined;

  // Helper function to extract block data whether it's a method or property
  const getBlockData = (block: any) => {
    return {
      meta: typeof block.meta === 'function' ? block.meta() : block.meta,
      owner: typeof block.owner === 'function' ? block.owner() : block.owner,
      num: typeof block.num === 'function' ? block.num() : block.num
    };
  };

  // Calculate the scale needed to fit the entire map in view
  const calculateFitScale = () => {
    if (!containerRef) return 1;
    const containerWidth = containerRef.clientWidth;
    const containerHeight = containerRef.clientHeight;
    const cellSize = 50;
    const gridGap = 1;
    const mapWidth = size.width * cellSize + (size.width - 1) * gridGap + 2; // cells + gaps + padding
    const mapHeight = size.height * cellSize + (size.height - 1) * gridGap + 2;
    
    // Add some padding around the map
    const paddingFactor = 0.85; // 85% of available space
    const scaleX = (containerWidth * paddingFactor) / mapWidth;
    const scaleY = (containerHeight * paddingFactor) / mapHeight;
    
    // Ensure minimum scale of 0.3 and maximum of 2
    return Math.max(0.3, Math.min(Math.min(scaleX, scaleY), 2));
  };

  // Fit the entire map to view
  const fitToView = () => {
    const fitScale = calculateFitScale();
    setScale(fitScale);
    setTranslateX(0);
    setTranslateY(0);
  };

  const handleWheel = (e: WheelEvent) => {
    e.preventDefault();
    const delta = e.deltaY > 0 ? 0.9 : 1.1;
    const newScale = Math.min(Math.max(0.3, scale() * delta), 3);
    setScale(newScale);
  };

  const handleMouseDown = (e: MouseEvent) => {
    setIsPanning(true);
    setStartX(e.clientX - translateX());
    setStartY(e.clientY - translateY());
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!isPanning()) return;
    setTranslateX(e.clientX - startX());
    setTranslateY(e.clientY - startY());
  };

  const handleMouseUp = () => {
    setIsPanning(false);
  };

  onMount(() => {
    if (containerRef) {
      containerRef.addEventListener('wheel', handleWheel, { passive: false });
      // Set initial scale to fit the entire map
      setTimeout(() => {
        fitToView();
      }, 100); // Small delay to ensure container dimensions are available
    }
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    
    // Expose fitToView function to parent
    if (onBoardRef) {
      onBoardRef({ fitToView });
    }
  });

  onCleanup(() => {
    if (containerRef) {
      containerRef.removeEventListener('wheel', handleWheel);
    }
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  });

  return (
    <div
      ref={containerRef}
      onMouseDown={handleMouseDown}
      class="overflow-hidden w-full flex-1 relative select-none"
      style={{
        background: '#2a2a2a',
        cursor: isPanning() ? 'grabbing' : 'grab',
      }}
    >
      <div
        style={{
          transform: `translate(${translateX()}px, ${translateY()}px) scale(${scale()})`,
          'transform-origin': '0 0',
          transition: isPanning() ? 'none' : 'transform 0.1s ease-out',
          position: 'absolute',
          top: '50%',
          left: '50%',
          'margin-top': `-${(size.height * 50) / 2}px`,
          'margin-left': `-${(size.width * 50) / 2}px`,
        }}
      >
        <div
          style={{
            display: 'grid',
            'grid-template-columns': `repeat(${size.width}, 50px)`,
            'grid-template-rows': `repeat(${size.height}, 50px)`,
            gap: '1px',
            background: '#1a1a1a',
            padding: '1px',
          }}
        >
          <For each={blocks ? blocks : Array.from({length: size.height}).map(() => Array(size.width).fill(null).map(() => createMockBlock('blank', 0, 0)))}>
            {(row, y) => (
              <For each={row}>
                {(block, x) => {
                  const blockData = getBlockData(block);
                  return (
                    <Cell 
                      meta={blockData.meta} 
                      owner={blockData.owner} 
                      num={blockData.num} 
                      position={{ x: x(), y: y() }}
                      isSelected={selectedCell()?.x === x() && selectedCell()?.y === y()}
                      onSelect={() => {
                        const selected = selectedCell();
                        if (selected && selected.x === x() && selected.y === y()) {
                          setSelectedCell(null);
                        } else {
                          setSelectedCell({ x: x(), y: y() });
                        }
                      }}
                    />
                  );
                }}
              </For>
            )}
          </For>
        </div>
      </div>
    </div>
  );
}

export default Board;
