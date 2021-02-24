import React from 'react';

class Cart extends React.Component {

  constructor (props) {
    super (props);

    this.state = {
      cart: null
    }
  }

  render() {
    return (
      <div>
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
        Add
      </button>
    );
  }
}


export default Cart;
