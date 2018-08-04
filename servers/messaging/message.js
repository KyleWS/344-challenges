"use strict";

const mongodb = require("mongodb")

class Message {
   constructor(messageData) {
      messageData = messageData || {};
      this._id = new mongodb.ObjectID();
      this.channelID = messageData.channelID;
      this.body = messageData.body;
      this.createdAt = messageData.createdAt;
      this.creator = messageData.creator;
      this.editedAt = messageData.editedAt;
   }

   // Weird bug: I need to add a empty string value to the
   // _id property, otherwise it is not stored properly
   // and cannot be found.. W/e if it works it works i guess
   getJSON() {
      return {
         _id: this._id + "",
         channelID: this.channelID,
         body: this.body,
         createdAt: this.createdAt,
         creator: this.creator,
         editedAt: this.editedAt
      }
   }
}

module.exports = Message;
