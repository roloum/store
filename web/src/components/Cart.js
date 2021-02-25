import React from 'react';

class Cart extends React.Component {

  constructor (props) {
    super (props)

    //bind button to App.js so we can display the list of items
    this.handleAddItemClick = this.handleAddItemClick.bind(props.parent);
    this.addItemOnClick = props.addItemOnClick;

    this.handleDeleteItemClick = this.handleDeleteItemClick.bind(props.parent);
    this.deleteItemOnClick = props.deleteItemOnClick;

    this.getCart = props.getCart


    this.state = {
      cart: props.cart
    }
  }

  handleAddItemClick () {

    this.addItemOnClick()

  }

  handleDeleteItemClick (item) {

    const cart = this.state.cart

    const endpoint = "/cart/"+cart.cart_id+"/items/"+item.item_id
    const url = "https://zqpjajqli1.execute-api.us-west-2.amazonaws.com/dev" + endpoint

    fetch(url, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then( (response) => {

      if (!response.ok) {
          throw new TypeError(`Error deleting item: `+ item.description);
      }

      return response.json()

    })
    .then( (result) => {
      this.deleteItemOnClick(result)
    })
    .catch(e => {
      //Display error message
      console.log("Error: ", e)
    });

  }

  render() {

    const cart = this.state.cart

    if (!cart || !cart.items || cart.items.length === 0) {
      return (
        <div>
          <div>
            <AddItemButton onClick={() => this.handleAddItemClick()}/>
          </div>
          <div>
          There are no items
          </div>
        </div>
      );
    }

    return (
      <div>
        <div>
          <AddItemButton onClick={() => this.handleAddItemClick()}/>
        </div>
        <div>
        <h2>Items</h2>
        </div>
        <div>
          <ul className="CartUL">
            {cart.items.map((item) => {
              return (
                <li className="CartListRow" key={item.item_id} >
                  <div>
                  <span className="CartListDesc">{item.description}</span>
                  <span className="CartListPrice">${item.price}</span>
                  </div>
                  <div>
                    <div>
                      <span className="ItemListPrice">{item.quantity}</span>
                    </div>
                    <div>
                      <DeleteItemButton onClick={() => this.handleDeleteItemClick(item)} />
                    </div>
                  </div>
                </li>
              );
            })}
          </ul>
        </div>
      </div>
    )
  }
}

class AddItemButton extends React.Component {

  constructor(props) {
    super(props);
    this.onClick = props.onClick
  }

  render () {
    return (
      <button onClick={this.onClick}>
        Add Items
      </button>
    );
  }
}

class DeleteItemButton extends React.Component {

  constructor(props) {
    super(props);
    this.onClick = props.onClick
  }

  render () {
    return (
      <button onClick={this.onClick}>
        Delete
      </button>
    );
  }
}


export default Cart;
