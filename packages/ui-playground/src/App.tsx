import { Router, Route } from '@solidjs/router';
import Home from './pages/Home';
import TestBoard from './pages/TestBoard';

function App() {
  return (
    <Router>
      <Route path="/" component={Home} />
      <Route path="/test-board" component={TestBoard} />
    </Router>
  );
}

export default App
