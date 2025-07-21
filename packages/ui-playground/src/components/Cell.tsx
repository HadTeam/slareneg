import { createSignal } from 'solid-js';
import type { Meta, Owner, Num } from '@slareneg/shared-types';

// Block icons using Unicode outline characters
const blockIcons: Record<string, string> = {
  blank: '',
  castle: '🏛', // Greek building as castle
  king: '♔',   // King chess piece
  mountain: '⛰', // Mountain emoji
  soldier: '◆', // Diamond for soldier
};

// Block background colors
const blockColors: Record<string, string> = {
  blank: '#3a3a3a',
  castle: '#4a3a5a',
  king: '#5a4a3a',
  mountain: '#2a2a2a',
  soldier: '#3a4a5a',
};

// Owner color mappings
const ownerColors: Record<number, string> = {
  0: '#666',      // Neutral
  1: '#4169E1',   // Player 1 - Royal Blue
  2: '#DC143C',   // Player 2 - Crimson
  3: '#228B22',   // Player 3 - Forest Green
  4: '#FF8C00',   // Player 4 - Dark Orange
};

interface CellProps {
  meta: Meta;
  owner?: Owner;
  num?: Num;
  position?: { x: number; y: number };
  isSelected?: boolean;
  onSelect?: () => void;
}

function Cell(props: CellProps) {
  const { meta, owner = 0, num = 0, isSelected = false, onSelect } = props;
  const [isHovered, setIsHovered] = createSignal(false);
  
  // Check if this is an unknown/fog block
  const isUnknown = meta.name === 'unknown' || meta.name === 'fog';
  const icon = isUnknown ? '' : (blockIcons[meta.name] || '');
  
  // Use darker color for fog/unknown blocks
  const bgColor = isUnknown ? '#1a1a1a' : (owner > 0 ? ownerColors[owner] : blockColors[meta.name]);
  const textColor = owner > 0 ? '#fff' : '#888';
  
  const handleClick = (e: MouseEvent) => {
    e.stopPropagation();
    onSelect?.();
  };

  return (
    <div
      onClick={handleClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      style={{
        background: bgColor,
        display: 'flex',
        'flex-direction': 'column',
        'align-items': 'center',
        'justify-content': 'center',
        width: '50px',
        height: '50px',
        border: `2px solid ${isSelected ? '#fff' : isHovered() ? ownerColors[owner] || '#888' : 'transparent'}`,
        'box-sizing': 'border-box',
        'font-size': '20px',
        position: 'relative',
        cursor: 'pointer',
        transition: 'all 0.1s ease',
        transform: isHovered() ? 'scale(1.05)' : 'scale(1)',
        color: textColor,
      }}
    >
      {icon && (
        <div style={{
          'font-size': '24px',
          'line-height': '1',
        }}>
          {icon}
        </div>
      )}
      {num > 0 && meta.name !== 'mountain' && meta.name !== 'blank' && (
        <div style={{
          'font-size': '12px',
          'font-weight': 'bold',
          'margin-top': icon ? '2px' : '0',
        }}>
          {num}
        </div>
      )}
    </div>
  );
}

export default Cell;

