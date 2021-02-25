import React from 'react';

class ItemsList extends React.Component {

  constructor(props) {
    super(props);

    //bind button to App.jss so we can display the cart
    this.handleAddClick = this.handleAddClick.bind(props.parent);

    this.handleBackClick = this.handleBackClick.bind(props.parent);

    this.addOnClick = props.addOnClick;
    this.backOnClick = props.backOnClick;

    this.state = {
      items: null,
      cartId: props.cartId
    };

  }

  handleAddClick(item) {

    const data = {
      "item_id": item.item_id,
      "description": item.description,
      "quantity": 1,
      "price": item.price
    }

    let url = "https://zqpjajqli1.execute-api.us-west-2.amazonaws.com/dev/cart"
    if (this.state.cartId !== null) {
      url += "/"+this.state.cartId
    }

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

      this.addOnClick(result)

    })
    .catch(e => {
      //Display error message
      console.log("Error: ", e)
      return
    });

  }

  handleBackClick() {
    this.backOnClick()
  }

  componentDidMount() {
    const url = "https://zqpjajqli1.execute-api.us-west-2.amazonaws.com/dev/items"
    fetch(url)
      .then( (response) => {

        if (!response.ok ||
            !response.headers.get("Content-Type") === "application/json") {
            throw new TypeError(`Error loading the items list`);
        }

        return response.json()

      })
      .then((result) => {
        this.setState({items: result.items});
      })
      .catch(e => {
        //Display error message
        console.log("Error: ", e)
      })

  }

  componentWillUnmount() {
  }

  render() {

    if (!this.state.items || this.state.items.length === 0) {
      return (
        <div>
        There are no items
        </div>
      );
    }

    return (
      <div>
        <div>
        <h2>Items</h2>
        </div>
        <div>
          <ul className="ItemsUL">
            <li className="ItemListRow">
              <span className="ItemListDesc">Description</span>
              <span className="ItemListPrice">Price</span>
              <span className="ItemListBtn"></span>
            </li>
            {this.state.items.map((item) => {
              return (
                <li className="ItemListRow" key={item.item_id} >
                  <span className="ItemListDesc">{item.description}</span>
                  <span className="ItemListPrice">${item.price}</span>
                  <span className="ItemListBtn">
                    <AddButton onClick={() => this.handleAddClick(item)} />
                  </span>
                </li>
              );
            })}
          </ul>
        </div>
        <div>
          <BackButton onclick={() => this.handleBackClick()}/>
        </div>
      </div>
    );
  }
}



class BackButton extends React.Component {

  constructor(props) {
    super(props);
    this.onClick = props.onClick
  }

  render () {
    return (
      <button onClick={this.onClick}>
        Back
      </button>
    );
  }
}


class AddButton extends React.Component {

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

export default ItemsList;
