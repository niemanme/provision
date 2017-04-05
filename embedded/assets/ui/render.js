/* Digital Rebar: Provision */
/* Copyright 2016, RackN */
/* License: Apache v2 */
/* jshint esversion: 6 */

$.ajaxSetup({
  headers : {
    'Authorization' : 'Bearer 66PM4ttu5nSD5ecRi_l4muqjpY7fJT3zSAMjJQ2o5E2h8sWoK4dXZWSOdfP3x53QMvTU6s_iuMwGFPtrLWEdFp03_qpsimoYcfrFoeOOx1tsLEzijNXMrkei3OvGggiSHF6US8nVq_l-A0TD83Y71jVwln7Rb7XCxy2mxHlEtgj8CFCEFYFPJlQ6Pp3EsoPOq7zo9rxCkZZp1VijN9b9xBLVedf-NYoBzJY1d9fHp5DDV-U8M70F-bZT3D5owfyikOniHf6j6OgUbGzklH_ZQ1wZl8GGbmyIltuTS2wDZxhmMB58PmCeJNNQ4UiIbdSez1RuA0UsMOD2YEpXWHDL6_w57VrQxu0zRLO5'
  }
});

class RSAuth extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      token: "66PM4ttu5nSD5ecRi_l4muqjpY7fJT3zSAMjJQ2o5E2h8sWoK4dXZWSOdfP3x53QMvTU6s_iuMwGFPtrLWEdFp03_qpsimoYcfrFoeOOx1tsLEzijNXMrkei3OvGggiSHF6US8nVq_l-A0TD83Y71jVwln7Rb7XCxy2mxHlEtgj8CFCEFYFPJlQ6Pp3EsoPOq7zo9rxCkZZp1VijN9b9xBLVedf-NYoBzJY1d9fHp5DDV-U8M70F-bZT3D5owfyikOniHf6j6OgUbGzklH_ZQ1wZl8GGbmyIltuTS2wDZxhmMB58PmCeJNNQ4UiIbdSez1RuA0UsMOD2YEpXWHDL6_w57VrQxu0zRLO5"
    };

    this.changeToken = this.changeToken.bind(this);
  }
  
  // get the subnets and interfaces once this component mounts
  componentDidMount() {
  }

  // makes the post/put request to update the subnet
  // also updates the interface
  changeToken() {
    $.ajaxSetup({
      headers : {
        'Authorization' : 'Bearer ${this.state.token}'
      }
    });
    var subnet = this.state.subnets[i];
  }

  render() {
    <input 
      type="password"
      name="token"
      size="10"
      value={this.state.token}
      onChange={this.handleChange} />
  }

}

ReactDOM.render(<RSAuth />, rsauth);
