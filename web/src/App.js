import React from 'react';

import './App.css'
import Cart from './components/Cart';
import ItemsList from './components/ItemsList';

class App extends React.Component {

  constructor (props) {
    super(props);

    this.state = {
      showCart: true,
      cartId: null,
      cart: null,
      subtotal: 0,
      itemCount: 0
    };

  }

  addItemOnClick () {
    this.setState({showCart: false});
  }

  backOnClick () {
    this.setState({showCart: true});
  }

  addOnClick (cart) {

    this.setState({
      showCart: true,
      cart: cart,
      cartId: cart.cart_id,
      subtotal: cart.total,
      itemCount: cart.count
    });

  }

  deleteItemOnClick (cart) {

    //Need to look find a better way to render the cart component
    //Updating the cart in the component state did not do the work
    //Unless the showCart flag is changed to false and then back to true
    //Doing this will cause an extra request of the items list
    this.setState({
      showCart: false,
      cart: cart,
      cartId: cart.cart_id,
      subtotal: cart.total,
      itemCount: cart.count
    });
    this.setState({
      showCart: true
    });
  }

  render() {

    const showCart = this.state.showCart;
    let section;

    if (showCart) {
      let props = {
        addItemOnClick: this.addItemOnClick,
        cart: this.state.cart,
        parent: this
      };
      section = <Cart {...props}/>;
    } else {

      let props = {
        addOnClick: this.addOnClick,
        backOnClick: this.backOnClick,
        cartId: this.state.cartId,
        parent: this
      };
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
