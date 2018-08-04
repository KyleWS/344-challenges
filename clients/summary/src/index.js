import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
//
// var ajax = new XMLHttpRequest();
// ajax.open("GET", "http://localhost:3000/v1/summary?url=" + input.value)
// ajax.onload = function(response) {
//    console.log(response.target.response)
// }
// ajax.onerror = function() {
//    console.log("error")
// }
// ajax.send();
