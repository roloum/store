import React from 'react';

class App extends React.Component {
  triggerFoo() {
    this.foo.toggle();
  }
  render() {
    return (
      <div>
        <Foo ref={foo => this.foo = foo} />
        <Button onClick={this.triggerFoo.bind(this)}/>
      </div>
    );
  }
}

class Foo extends React.Component {
  state = {foo: false}
  toggle() {
    this.setState({
      foo: !this.state.foo
    });
  }
  render() {
    return (
      <div>
        Foo Triggered: {this.state.foo.toString()}
      </div>
    );
  }
}


class Button extends React.Component {
  render() {
    return (
      <button onClick={this.props.onClick}>
        Click This
      </button>
    );
  };
}

export default App;
