import React from 'react';

import './App.css'
import Header from './components/Header';
import Cart from './components/Cart';
import ItemsList from './components/ItemsList';

class App extends React.Component {

  constructor (props) {
    super(props)

    this.state = {
      showCart: false,
      cartId: null
    };

  }
  addItemOnClick (itemId) {
    console.log("onAddItemClick sent from App id: ", itemId)
  }

  render() {
    const showCart = this.state.showCart;
    let section;

    if (showCart) {
      section = <Cart />;
    } else {
      section = <ItemsList addItemOnClick={this.addItemOnClick}/>;
    }

    return (
      <div className="App">
        <Header />
        {section}
      </div>
    );
  }

}

export default App;
