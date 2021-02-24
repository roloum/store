import React from 'react';

class Header extends React.Component {
  render() {
    return (
      <div className="Header">
        <div className="HeaderTitle">
          <h1>Shopping Cart</h1>
        </div>
        <div className="CartItemsCount">
          <h2>Items: 0</h2>
        </div>
      </div>
    );
  }
}

export default Header;
