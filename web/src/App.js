import React from 'react';

import './App.css'
import Cart from './components/Cart';
import ItemsList from './components/ItemsList';

class App extends React.Component {

  constructor (props) {
    super(props)

    this.state = {
      showCart: true,
      cartId: null,
      cart: null,
      subtotal: 0,
      itemCount: 0
    };

  }

  addItemOnClick () {
    console.log("add item on click app.js hide cart")
    this.setState({showCart: false})
  }

  backOnClick () {
    console.log("back on click app.js show cart")
    this.setState({showCart: true})
  }

  addOnClick (cart) {

    this.setState({
      showCart: true,
      cart: cart,
      cartId: cart.cart_id,
      subtotal: cart.total,
      itemCount: cart.count
    })
  }

  render() {

    const showCart = this.state.showCart;
    let section;

    if (showCart) {
      let props = {
        addItemOnClick: this.addItemOnClick,
        cart: this.state.cart,
        parent: this
      }
      section = <Cart {...props}/>;
    } else {

      let props = {
        addOnClick: this.addOnClick,
        backOnClick: this.backOnClick,
        cartId: this.state.cartId,
        parent: this
      }
      section = <ItemsList {...props}/>;
    }

    return (
      <div className="App">


        <div className="Header">
          <div className="HeaderTitle">
            <h1>Shopping Cart</h1>
          </div>
          <div className="CartItemsCount">
          <div>
            <b>Count: {this.state.itemCount}</b>
          </div>
          <div>
            <b>Subtotal: ${this.state.subtotal}</b>
          </div>
          </div>
        </div>


        {section}
      </div>
    );
  }

}

export default App;
