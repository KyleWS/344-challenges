import React, { Component } from 'react';
import './App.css';
import $ from 'jquery';

class App extends Component {
   constructor(props) {
      super(props);
      this.state = {
         data: ""
      };
      this.handleData = this.handleData.bind(this);
   }

   handleData(data) {
      this.setState({
         data: JSON.stringify(data)
      });
   }

  render() {
    return (
      <div id="app-body">
         <Toolbar callback={this.handleData} />
         <Board data={this.state.data} />
      </div>
    );
  }
}

class Board extends Component {
   render() {
      if(this.props.data !== "") {
         var obj = JSON.parse(this.props.data);
         return (
            <div>
               <Card source={obj} />
            </div>
         );
      } else {
         return (
            <h1>Enter a URL</h1>
         );
      }
   }
}

class Card extends Component {
   render() {
      if (this.props.source.images.length > 0) {
         var imgs = this.props.source.images.forEach(function(img) {
         <img src={img.url} alt="" />});
      }
      return (
         <div className="card-body">
            <h1>{this.props.source.title}</h1>
            <h2>Type: {this.props.source.type}</h2>
            <h2>url: {this.props.source.url}</h2>
            <h3>Description</h3>
            <p>{this.props.source.description}</p>
            <h3>Images</h3>
            <p>{imgs}</p>


         </div>
      )
   }
}

class Toolbar extends Component {
   checkURL(address) {
      if(address.startsWith("https://")) {
         address = address;
      } else if(!address.startsWith("http://")) {
         address = "http://" + address;
      }
      return address;
   }

   makeHTTPRequestCallback(address, callback) {
      if(address !== "") {
         address = this.checkURL(address);
         $.ajax({url: "/v1/summary?url=" + address , success: function(data) {
            callback(data)
         }})
      }
   }

   render() {
      return (
         <div className="toolbar-body">
            <input id="toolbar-search" type="text" placeholder="enter url and press go"/>
            <button id="toolbar-go-fetch"
               onClick={(e) => this.makeHTTPRequestCallback(
                  document.getElementById("toolbar-search").value,
                  this.props.callback
               )}>Go</button>
         </div>
      )
   }
}

export default App;
