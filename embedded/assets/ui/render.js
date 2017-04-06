/* Digital Rebar: Provision */
/* Copyright 2016, RackN */
/* License: Apache v2 */
/* jshint esversion: 6 */

class Token extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      token: "",
      access: false
    };

    this.handleChange = this.handleChange.bind(this);
  }
  
  // get the page and interfaces once this component mounts
  componentDidMount() {
    console.debug("foo");
  }

  // called when an input changes
  handleChange(event) {
    this.state.token = event.target.value;
    $.ajaxSetup({
      headers : {
        'Authorization' : 'Bearer ' + this.state.token
      }
    });
    // get the interfaces from the api
    console.debug("barbar");
    $.getJSON("../api/v3/prefs", data => { 
      console.debug(data);
      this.setState({access: true});
    }).fail((jqXHR, textStatus, errorThrown) => {
      console.debug("failed")
      this.setState({access: false});
    });
    // get the interfaces from the api
    console.debug(this.state.token);
  }

  render() {
    return (
      <div>
      <input
        type="text"
        name="token"
        size="15"
        placeholder="no_token"
        value={this.state.token}
        onChange={this.handleChange} />
        <div>
        {this.state.access ? "yes" : 'no'}
        </div>
        </div>
    );
  }
}

class Page extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      access: false
    };

    this.handleChange = this.handleChange.bind(this);
  }
  
  // get the page and interfaces once this component mounts
  componentDidMount() {
    console.debug("bar");
  }

  // called when an input changes
  handleChange(event) {
    console.debug("barbar");
  }

  render() {
    return (
      <div>
        {this.props.access 
          ? "Access Granted"
          :  <center>Access Token: <Token access={this.props.access} /></center>
        }
        <div> {this.state.access ? "YES" : 'NO'}
        </div>
      </div>
    );
  }
}

ReactDOM.render(<Page />, page);