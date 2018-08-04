"use strict";

var CONNECT_CONST = "api.kylewilliamscreates.com"

window.onload = function() {
   if (localStorage.getItem("Authorization") == null) {
      window.location = "index.html";
   } else {
      console.log(localStorage.getItem("Authorization"));
   }

   // sign out button
   $("#sign-out").click(function() {
      $.ajax({
         url: "https://" + CONNECT_CONST + "/v1/sessions/mine",
         type: "DELETE",
         headers: {"Authorization": localStorage.getItem("Authorization")},
         success: function(data) {
            console.log(data);
         },
         error: function(data) {
            errorFill(data.responseText);
         }
      }).then(function(data, status, xhr) {
         localStorage.removeItem("Authorization");
         window.location = "index.html";
      });
   });

   $.ajax({
      url: "https://" + CONNECT_CONST + "/v1/users/me",
      type: "GET",
      headers: {"Authorization": localStorage.getItem("Authorization")},
      success: function(data) {

      },
      error: function(data) {
         errorFill(data.responseText);
      }
   }).then(function(data, status, xhr) {

      var userInfo = JSON.parse(data);
      $("#username").html(userInfo.userName);
      $("#firstname").html(userInfo.firstName);
      $("#lastname").html(userInfo.lastName);
   });

   $("#update").click(function() {
      window.location = "update.html";
   });

   // sign out button
   $("#search").click(function() {
      var param = $("#search-query").val();
      $.ajax({
         url: "https://" + CONNECT_CONST + "/v1/users?q=" + param,
         type: "GET",
         headers: {"Authorization": localStorage.getItem("Authorization")},
         success: function(data) {

         },
         error: function(data) {
            errorFill(data.responseText);
         }
      }).then(function(data, status, xhr) {
         errorFill("");
         data = JSON.parse(data)
         $("#search-results").html("");
         for (var i = 0; i < data.length; i++) {
            console.log(entry)
            var entry = $("<div>", {"class": "result"});
            entry.html(data[i].email + " " + data[i].userName + " " + data[i].lastName);
            $("#search-results").append(entry);
         }
      });
   });
}

function errorFill(errMsg) {
   $("#error-div").html(errMsg);
}
