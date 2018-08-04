"use strict";

var CONNECT_CONST = "api.kylewilliamscreates.com"

window.onload = function() {
      $('#test').click(function() {
         $.post("https://"+CONNECT_CONST+"/v1/users", function() {
            console.log("sent request I guess");
         }).done(function() {
            console.log("success");
         }).fail(function(data) {
            errorFill(data.responseText);
         });
      });
      // NEW ACCOUNT
      $('#new-account-submit').click(function(e) {
         var myJson = {
            email: $("#new-account-email").val(),
            password: $("#new-account-password").val(),
            passwordConf: $("#new-account-passwordConf").val(),
            userName: $("#new-account-userName").val(),
            firstName: $("#new-account-firstName").val(),
            lastName: $("#new-account-lastName").val()
         }
         $.ajax({
            url: "https://"+CONNECT_CONST+"/v1/users",
            dataType: "json",
            contentType: "application/json; charset=UTF-8",
            data: JSON.stringify(myJson),
            type: "POST",
            headers: {"Authorization": localStorage.getItem("Authorization")},
            success: function(data) {
               console.log(data);
            },
            error: function(data) {
               errorFill(data.responseText);
            }
         }).then(function(data, status, xhr) {
            successFill(data.responseText);
            console.log(xhr.getResponseHeader("authorization"));
            localStorage.setItem("Authorization", xhr.getResponseHeader("Authorization"));
            window.location = "signed-in.html";
         });
      });
      // SIGN IN TO EXISTING ACCOUNT
      $('#existing-account-submit').click(function(e) {
         var myJson = {
            email: $("#existing-account-email").val(),
            password: $("#existing-account-password").val()
         }
         $.ajax({
            url: "https://"+CONNECT_CONST+"/v1/sessions",
            dataType: "json",
            contentType: "application/json; charset=UTF-8",
            data: JSON.stringify(myJson),
            type: "POST",
            headers: {"Authorization": localStorage.getItem("Authorization")},
            success: function(data) {
               console.log(data);
            },
            error: function(data) {
               errorFill(data.responseText);
            }
         }).then(function(data, status, xhr) {
            successFill(data.responseText);
            console.log(xhr.getResponseHeader("Authorization"));
            localStorage.setItem("Authorization", xhr.getResponseHeader("Authorization"));
            window.location = "signed-in.html";
         });
      });
}

function errorFill(errMsg) {
   $("#error-div").html(errMsg);
}

function successFill(successMsg) {
   $("#success-div").html(successMsg);
}
