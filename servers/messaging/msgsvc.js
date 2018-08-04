#! /usr/bin/env node
"use strict";

const XUSER = "xuser placeholder";

const ErrChannelExists = "error channel name already being used";
const ErrChannelNotFound = "error could not find channel";
const ErrMessageNotFound = "error could not find message";
const ErrNameRequired = "error name field required";
const ErrChannelRequired = "error channelID field required";
const ErrBodyRequired = "error body field required";
const ErrStartCheckChannels = "error getting default general channel";

const StatusXUserRequired = "authenticated X-User field not present in header";
const StatusAuthorRequired = "error must be creator of this entity to alter/delete it";

const ContentType = "Content-type";
const applicationJson = "application/json";
const XUser = "X-User";

// libs
const express = require("express");
const morgan = require("morgan");
const mongodb = require("mongodb");

// models
const Channel = require("./channel.js");
const Message = require("./message.js");
const MsgStore = require("./msgstore.js");

const msgAddr = process.env.MSGADDR || "localhost:5000";
const [host, port] = msgAddr.split(":");

const mongoAddr = process.env.DBADDR || "localhost:27017";
const mongoDbName = "MSGSVCDB"
const mongoURL = `mongodb://${mongoAddr}/${mongoDbName}`

const app = express();
app.use(morgan("dev"));

// Remember to reject unauthenticated persons
mongodb.MongoClient.connect(mongoURL)
   .then(db => {
      let msgStore = new MsgStore(db, "messagesCollection", "channelsCollection");

      // parse json from input
      app.use(express.json());

      app.get("/v1/channels", (req, res) => {
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            msgStore.getChannels({})
               .then(results => {
                  res.json(results);
               }).catch(err => {
                  res.json({error: err.message});
               });
         }
      });

      app.post("/v1/channels", (req, res) => {
         // create new channel using json body of req. name required, description not.
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            if (!req.body.name) {
               res.json({error: ErrNameRequired});
            } else {
               msgStore.getChannels({"name": req.body.name}).then(results => {
                  if (results.length == 0) {
                     let newChannel = new Channel({
                        name: req.body.name,
                        description: req.body.description || "",
                        createdAt: Date.now(),
                        creator: req.get(XUser), // THIS SHIT RIGHT HERE
                        editedAt: Date.now()
                     });
                     msgStore.insertChannel(newChannel)
                        .then(channel => {
                           res.json(channel);
                        })
                        .catch(err => {
                           res.json({error: err.message});
                        });
                  } else {
                     res.json({error: ErrChannelExists});
                  }
               }).catch(err => {
                  res.json({error: err.message});
               });
            }
         }
      });

      app.get("/v1/channels/:id", (req, res) => {
         // respond with latest 50 messages from given channel id, in json array
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            msgStore.getMessages({"channelID": req.params.id})
               .then(results => {
                  res.json(results);
               }).catch(err => {
                  res.json({error: err.message});
               });
         }
      });

      app.post("/v1/channels/:id", (req, res) => {
         // create new message in the provided channel.
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            if (!req.body.body) {
               res.json({error: ErrBodyRequired});
            } else {
               // check to make sure this channel exists
               msgStore.getChannels({"_id": req.params.id}).then(results => {
                  // Presumably we have one result which is the channel we want to post to.
                  if (results.length > 0) {
                     let newMessage = new Message({
                        channelID: req.params.id,
                        body: req.body.body,
                        createdAt: Date.now(),
                        creator: req.get(XUser),
                        editedAt: Date.now()
                     });
                     msgStore.insertMessage(newMessage)
                           .then(message => {
                              res.json(message);
                           })
                           .catch(err => {
                              res.json({error: err.message});
                           });
                  } else {
                     res.json({error: ErrChannelNotFound})
                  }
               })
               .catch(err => {
                  res.json({error: err.message});
               });
            }
            // respond with copy of this new messages
            // on property you should read from request is "body"
         }
      });

      app.patch("/v1/channels/:id", (req, res) => {
         // if the current user created the channel, update only the name/description
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            msgStore.getChannels({"_id": req.params.id}).then(results => {
               let userJSON = JSON.parse(results[0].creator);
               let xuserJSON = JSON.parse(req.get(XUser));
               if (userJSON.id != xuserJSON.id) {
                     res.status(403).send(StatusAuthorRequired)
               } else {
                  if (results.length > 0) {
                     let updates = {
                        name: req.body.name || results[0].name,
                        description: req.body.description || results[0].description,
                        editedAt: Date.now()
                     };
                     msgStore.updateChannel(req.params.id, updates)
                        .then(result => {
                           res.json(result);
                        })
                        .catch(err => {
                           res.json({error: err.message});
                        })
                  } else {
                     res.json({error: ErrChannelNotFound});
                  }
               }
            })
            .catch(err => {
               res.json({error: err.message});
            });
            // fields using the json provided. Respond with copy of newly updated channel
            // as json. If the current user isnt the channel creator, throw dat 403
         }
      });

      app.delete("/v1/channels/:id", (req, res) => {
         // If current user created channel, delete it and all related messages.
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set({ContentType, applicationJson});
            msgStore.getChannels({"_id": req.params.id}).then(results => {
               let userJSON = JSON.parse(results[0].creator);
               let xuserJSON = JSON.parse(req.get(XUser));
               if (userJSON.id != xuserJSON.id) {
                     res.status(403).send(StatusAuthorRequired);
               } else {
                  msgStore.deleteChannel(req.params.id).then(result => {
                     msgStore.deleteMessages({"channelID": req.params.id})
                        .then(results => {
                           res.json({results: results, result: result});
                        })
                        .catch(err => {
                           res.json({result: result, error: err.message});
                        });
                  })
                  .catch(err => {
                     res.json({error: err.message});
                  });
                  // Otherwise, respond with 403
               }
            }).catch(err => {
               res.json({error: ErrChannelNotFound})
            });
            }
      });

      app.patch("/v1/messages/:id", (req, res) => {
         // if the current user created the message, update the message body
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set(ContentType, applicationJson);
            msgStore.getMessages({"_id": req.params.id}).then(results => {
               let userJSON = JSON.parse(results[0].creator);
               let xuserJSON = JSON.parse(req.get(XUser));
               console.log("we make it this far");
               if (userJSON.id != xuserJSON.id) {
                  console.log("%s and %s", userJSON.id, xuserJSON.id);
                     res.status(403).send(StatusAuthorRequired);
               } else {
                  if (results.length > 0) {
                     let updates = {
                        body: req.body.body || results[0].body,
                        editedAt: Date.now()
                     };
                     msgStore.updateMessage(req.params.id, updates)
                        .then(result => {
                           res.json(result);
                        })
                        .catch(err => {
                           res.json({error: err.message});
                        })
                  } else {
                     res.json({error: ErrMessageNotFound});
                  }
               }
            })
            .catch(err => {
               res.json({error: err.message});
            });
            // and respond with a copy.  Also respond with 403 is hax
         }
      });

      app.delete("/v1/messages/:id", (req, res) => {
         if (req.get(XUser) == undefined) {
            res.status(401).send(StatusXUserRequired);
         } else {
            res.set({ContentType, applicationJson});
            msgStore.getMessages({"_id": req.params.id}).then(results => {
               let userJSON = JSON.parse(results[0].creator);
               let xuserJSON = JSON.parse(req.get(XUser));
               if (userJSON.id != xuserJSON.id) {
                     res.status(403).send(StatusAuthorRequired);
               } else {
                  msgStore.deleteMessages({"_id": req.params.id})
                     .then(results => {
                        if (results.n > 0) {
                           res.json({error: ErrMessageNotFound});
                        } else {
                           res.json({results: results});
                        }
                     })
                     .catch(err => {
                        res.json({error: err.message});
                     });
               }
            }).catch(err => {
               res.json({error: err.message});
            });

         }
      });

      app.listen(port, host, () => {
         msgStore.getChannels({"name": "general"}).then(result => {
            if (result.length == 0) {
               let generalChannel = new Channel({
                  name: "general",
                  description: "generic channel",
                  createdAt: Date.now(),
                  creator: "server", // THIS SHIT RIGHT HERE
                  editedAt: Date.now()
               });
               msgStore.insertChannel(generalChannel).then(result =>{})
                  .catch(err => {throw err;});
            }
         })
         .catch(err => {
            throw err;
         });
         console.log(`message server listening on http://${msgAddr}`);
      });
   })
   .catch(err => {
      throw err;
   });
