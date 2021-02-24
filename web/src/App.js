import React from 'react';

import './App.css'
import Header from './components/Header';
import Cart from './components/Cart';
import ItemsList from './components/ItemsList';

class App extends React.Component {

  render() {
    return (
      <div className="App">
        <Header />
        <ItemsList />
      </div>
    );
  }

}

export default App;
