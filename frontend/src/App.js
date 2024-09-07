import React from 'react';
import './App.css';
import StarknetAbiFetcher from './StarknetAbiFetcher';
import StarknetBackfillFetcher from './StarknetBackfillFetcher';
function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>StarkNet ABI Fetcher</h1>
      </header>
      <main>
        <StarknetAbiFetcher />
        <StarknetBackfillFetcher/>
      </main>
    </div>
  );
}

export default App;
