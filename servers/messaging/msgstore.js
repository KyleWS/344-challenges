"use strict";

const mongodb = require("mongodb");

class MsgStore {
    constructor(db, messageCollectionName, channelCollectionName) {
        this.msgCol = db.collection(messageCollectionName);
        this.chanCol = db.collection(channelCollectionName);
    }

   insertChannel(channel) {
      return this.chanCol.insertOne(channel.getJSON())
         .then(() => channel);
   }

   getChannels(query) {
      return this.chanCol.find(query)
         .toArray();
   }

   getMessages(query) {
      return this.msgCol.find(query)
         .limit(50).toArray();
   }

   insertMessage(message) {
      return this.msgCol.insertOne(message.getJSON())
         .then(() => message);
   }

   updateChannel(id, updates) {
      let updateDoc = {
         "$set": updates
      }
      return this.chanCol.findOneAndUpdate(
         {_id: id},
         updateDoc,
         {returnOriginal: false})
         .then(result => result.value);
   }

   updateMessage(id, updates) {
      let updateDoc = {
         "$set": updates
      }
      return this.msgCol.findOneAndUpdate(
         {_id: id},
         updateDoc,
         {returnOriginal: false})
         .then(result => result.value);
   }

   deleteChannel(id) {
      return this.chanCol.deleteOne({_id: id})
         .then(result => result);
   }

   deleteMessages(query) {
      return this.msgCol.deleteMany(query)
         .then(results => results);
   }
}

module.exports = MsgStore;
