import { createSignal } from 'solid-js';
import { A } from '@solidjs/router';
import Board from '../components/Board';
import { generateMockMap } from '../mocks/mockMap';
import type { Size } from '@slareneg/shared-types';

function TestBoard() {
  const boardSize: Size = { width: 20, height: 20 };
  const [blocks] = createSignal(generateMockMap());

  return (
    <div style={{ 
      height: '100vh',
      display: 'flex',
      'flex-direction': 'column',
      margin: 0,
      padding: 0,
      overflow: 'hidden'
    }}>
      <div style={{ 
        padding: '10px',
        background: '#f0f0f0',
        'border-bottom': '1px solid #ddd',
        display: 'flex',
        'align-items': 'center',
        gap: '20px'
      }}>
        <A href="/" style={{
          padding: '5px 15px',
          background: '#666',
          color: 'white',
          'text-decoration': 'none',
          'border-radius': '3px',
          'font-size': '14px'
        }}>
          ← Back to Home
        </A>
        <div>
          <h2 style={{ margin: '0 0 5px 0' }}>Test Board - Slareneg</h2>
          <p style={{ margin: 0, 'font-size': '14px' }}>Click to select • Drag to move • Scroll to zoom</p>
        </div>
      </div>
      <Board blocks={blocks()} size={boardSize} />
    </div>
  );
}

export default TestBoard;
