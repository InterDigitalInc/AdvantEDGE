import React from 'react';

export default class ResizeableContainer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      width: 500,
      height: 500
    };

    this.sizingCallback = element => {
      if (element) {
        const width = element.getBoundingClientRect().width;
        const height = element.getBoundingClientRect().height;
        this.setState({
          height: height,
          width: width
        });
      }
    };
  }

  render() {
    return (
          <>
            <div      
              ref={this.sizingCallback}
            >
              {this.props.children(this.state.width, this.state.height)}   
            </div>
          </>       
    );
  }
}