import { A } from '@solidjs/router';

function Home() {
  return (
    <div style={{
      height: '100vh',
      display: 'flex',
      'flex-direction': 'column',
      'align-items': 'center',
      'justify-content': 'center',
      padding: '20px'
    }}>
      <h1 style={{ 'margin-bottom': '20px' }}>Welcome to Slareneg</h1>
      <nav>
        <A href="/test-board" style={{
          display: 'inline-block',
          padding: '10px 20px',
          background: '#0066cc',
          color: 'white',
          'text-decoration': 'none',
          'border-radius': '5px',
          'font-size': '16px'
        }}>
          Go to Test Board
        </A>
      </nav>
    </div>
  );
}

export default Home;
