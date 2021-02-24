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
      cartId: null,
      cart: null,
    };

  }
  addItemOnClick (item) {

    const data = {
      "item_id": item.item_id,
      "description": item.description,
      "quantity": 1,
      "price": item.price
    }

    console.log("Adding item from App: ", data)

    const url = "https://zqpjajqli1.execute-api.us-west-2.amazonaws.com/dev/cart"

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
