import React from 'react';

class Cart extends React.Component {

  constructor (props) {
    super (props)

    //bind button to App.js so we can display the list of items
    this.handleAddItemClick = this.handleAddItemClick.bind(props.parent);

    this.addItemOnClick = props.addItemOnClick;

    console.log("cart.js displaying cart on props")
    console.log(props.cart)

    this.state = {
      cart: props.cart
    }
    console.log("cart.js displaying cart on state")
    console.log(this.state.cart)
  }

  handleAddItemClick () {

    this.addItemOnClick()

  }

  render() {
    console.log("cart.js displaying items on render")
    console.log(this.state.cart)

    if (!this.state.cart || !this.state.cart.items || this.state.cart.items.length === 0) {
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
          <ul className="ItemsUL">
            {this.state.cart.items.map((item) => {
              return (
                <li className="ItemListRow" key={item.item_id} >
                  <span className="ItemListDesc">{item.description}</span>
                  <span className="ItemListPrice">{item.quantity}</span>
                  <span className="ItemListPrice">${item.price}</span>
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


export default Cart;
