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

    console.log("cartId on items list ", this.state.cartId)
  }

  handleAddClick(item) {
    this.addOnClick(item)
  }

  handleBackClick() {
    console.log("handleBackClick on list.js")
    this.backOnClick()
  }

  componentDidMount() {
    /*
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
      */
    this.setState({items: [{"item_id":"0dbe71c6-8584-43cd-be13-69ddf5651289","description":"Water bottle","price":4},{"item_id":"5408ea4e-1674-484a-947c-721e205b7d7f","description":"Book","price":10.99},{"item_id":"609544d0-1d17-4739-8056-9432bfd197bc","description":"Umbrella","price":7.29},{"item_id":"83adae8c-adee-4729-974d-452c8c30aa6c","description":"Pen","price":0.99},{"item_id":"9008e368-b2e0-4fe6-a677-33148a4af036","description":"Headphones","price":17.99},{"item_id":"b448e2a1-abd0-4a92-80e3-523fc0929487","description":"Jacket","price":59.99}]})

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
