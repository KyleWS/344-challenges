"use strict";

var CONNECT_CONST = "api.kylewilliamscreates.com"

window.onload = function() {
   if (localStorage.getItem("Authorization") == null) {
      window.location = "index.html";
   } else {
      console.log(localStorage.getItem("Authorization"));
   }

   $("#back").click(function() {
      window.location = "signed-in.html";
   });

   // Update Profile
   $('#submit').click(function(e) {
      var myJson = {
         firstName: $("#firstname").val(),
         lastName: $("#lastname").val()
      }
      $.ajax({
         url: "https://"+CONNECT_CONST+"/v1/users/me",
         dataType: "json",
         contentType: "application/json; charset=UTF-8",
         data: JSON.stringify(myJson),
         type: "PATCH",
         headers: {"Authorization": localStorage.getItem("Authorization")},
         success: function(data) {
            console.log(data);
         },
         error: function(data) {
            errorFill(data.responseText);
         }
      }).then(function(data, status, xhr) {
         window.location = "signed-in.html";
      });
   });
}

function errorFill(errMsg) {
   $("#error-div").html(errMsg);
}
