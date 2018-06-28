# chatter

The purpose of the repository is to make a simple chat application using web sockets in Go.

### what i learned
1. How to use gorilla websocket library.

2. How to use gorilla cookies library.

3. How to do concurrency with Go using channels.

   Here is the main loop of the chat server. This method is an infinite loop which listens to all the channels in the Server struct.
   It runs without blocking by simply doing ```go server.Listen()```. 
   
  ```go
  // Listen ...
  // infinite loop listening to channels
  func (server *Server) Listen() {
    log.Println("Chat Server Listening .....")

    server.Router.HandleFunc("/chat", server.handleChat)
    server.Router.HandleFunc("/getAllMessages", server.handleGetAllMessages)

    for {
      select {
      // Adding a new user
      case user := <-server.addUser:
        log.Println("Added a new User")
        server.connectedUsers[user.id] = user
        log.Println("Now ", len(server.connectedUsers), " users are connected to chat room")

      case user := <-server.removeUser:
        log.Println("Removing user from chat room")
        delete(server.connectedUsers, user.id)

      case msg := <-server.newIncomingMessage:
        if len(server.Messages) > 5 {
          server.shiftMessages(msg)
        } else {
          server.Messages = append(server.Messages, msg)
        }
        server.sendAll(msg)
      case err := <-server.errorChannel:
        log.Println("Error : ", err)
      case <-server.doneCh:
        return
      }
    }
  }
  
  ```

4. Adding custom middleware to routes.

   Here is an example of adding a login requirment to certain routes

  
  ```go
  r.HandleFunc("/", reqLogin(indexHandler))           // Require Login
  r.HandleFunc("/login", loginHandler)
  r.HandleFunc("/getuser", reqLogin(getuserHandler))  // Require Login
  
  ...
  
  // This function recieves a handler and checks to see if the user is logged in.
  // If the use is not logged it they will be rerouted to "/login", otherwise
  // the handler the function recieves will be executed like normal. 
  func reqLogin(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      session, _ := store.Get(r, "auth")
      if auth, ok := session.Values["auth"].(bool); !ok || !auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
      }
      f(w, r)
    }
  }
  ```

### get it running

```
go get https://github.com/arman-ashrafian/chatter
go build
./chatter
```
