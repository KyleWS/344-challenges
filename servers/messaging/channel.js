"use strict";

const mongodb = require("mongodb");

class Channel {
   constructor(channelData) {
      channelData = channelData || {};
      this._id = new mongodb.ObjectID();
      this.name = channelData.name;
      this.description = channelData.description;
      this.createdAt = channelData.createdAt;
      this.creator = channelData.creator;
      this.editedAt = channelData.editedAt;
   }

   getJSON() {
      return {
         _id: this._id + "",
         name: this.name,
         description: this.description,
         createdAt: this.createdAt,
         creator: this.creator,
         editedAt: this.editedAt
      }
   }
}

module.exports = Channel;
