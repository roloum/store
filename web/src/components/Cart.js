import React from 'react';

class Cart extends React.Component {

  constructor (props) {
    super (props)

    //bind button to App.js so we can display the list of items
    this.handleAddItemClick = this.handleAddItemClick.bind(props.parent);

    this.addItemOnClick = props.addItemOnClick;

    this.state = {
      cart: props.cart
    }
  }

  handleAddItemClick () {

    this.addItemOnClick()

  }

  render() {
    return (
      <div>
        <AddItemButton onClick={() => this.handleAddItemClick()}/>
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
