import React from 'react';

import './App.css'
import Header from './components/Header';
import Cart from './components/Cart';
import ItemsList from './components/ItemsList';

class App extends React.Component {

  constructor (props) {
    super(props)

    this.state = {
      showCart: true,
      cartId: null,
      cart: null
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

  addOnClick (item) {

    const data = {
      "item_id": item.item_id,
      "description": item.description,
      "quantity": 1,
      "price": item.price
    }

    console.log("Adding item from App: ", data)

    let url = "https://zqpjajqli1.execute-api.us-west-2.amazonaws.com/dev/cart"
    if (this.state.cartId !== null) {
      url += "/"+this.state.cartId
    }

    console.log(url)

    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: JSON.stringify(data) // body data type must match "Content-Type" header
      })
    .then( (response) => {

      if (!response.status === 201 ||
          !response.headers.get("Content-Type") === "application/json") {
          throw new TypeError(`Error adding item: `+ data.description);
      }

      return response.json()

    })
    .then( (result) => {
      console.log("logging result")
      console.log(result)
      console.log("show cart true")
      this.setState({showCart: true, cart: result, cartId: result.cart_id})
      //this.forceUpdate()
    })
    .catch(e => {
      //Display error message
      console.log("Error: ", e)
      return
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
        <Header />
        {section}
      </div>
    );
  }

}

export default App;
