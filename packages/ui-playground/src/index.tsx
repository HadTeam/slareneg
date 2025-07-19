/* @refresh reload */
import { render } from 'solid-js/web'
import App from './App.tsx'
import './index.css'

// Global styles
const style = document.createElement('style');
style.textContent = `
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }
  
  html, body {
    width: 100%;
    height: 100%;
    overflow: hidden;
  }
  
  #root {
    width: 100%;
    height: 100%;
  }
`;
document.head.appendChild(style);

const root = document.getElementById('root')

render(() => <App />, root!)
