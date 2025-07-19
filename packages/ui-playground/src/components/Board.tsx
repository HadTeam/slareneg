import { For, createSignal, onMount, onCleanup } from 'solid-js';
import type { Blocks, Size } from '@slareneg/shared-types';
import Cell from './Cell';
import { createMockBlock } from '../mocks/mockBlocks';

interface BoardProps {
  blocks?: Blocks;
  size: Size;
}

function Board(props: BoardProps) {
  const { blocks, size } = props;
  const [scale, setScale] = createSignal(1);
  const [translateX, setTranslateX] = createSignal(0);
  const [translateY, setTranslateY] = createSignal(0);
  const [isPanning, setIsPanning] = createSignal(false);
  const [startX, setStartX] = createSignal(0);
  const [startY, setStartY] = createSignal(0);
  const [selectedCell, setSelectedCell] = createSignal<{x: number, y: number} | null>(null);

  let containerRef: HTMLDivElement | undefined;

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
    }
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
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
                {(block, x) => (
                  <Cell 
                    meta={block.meta()} 
                    owner={block.owner()} 
                    num={block.num()} 
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
                )}
              </For>
            )}
          </For>
        </div>
      </div>
    </div>
  );
}

export default Board;
