package controllers

import (
	"github.com/revel/revel"
  "golang.org/x/net/websocket"
  "ws_app/app/controllers/chatroom"
)


// import "fmt"


type IndexController struct {
	*revel.Controller
}


		/*
		|-----------------------------------------
		|  Indexpage
		|
		|---------------------------------------------
		*/
		func (this IndexController) Index() revel.Result {

        this.ViewArgs["title"] = "Index"
        // this.ViewArgs["user"] = "Nikeweke"
        return this.RenderTemplate("App/index.html")
		}



		/*
		|-----------------------------------------
		|  ChatRoom
		|
		|---------------------------------------------
		*/
		func (this IndexController) ChatRoom() revel.Result {

				this.ViewArgs["title"] = "WSockets"
				// this.ViewArgs["user"] = "Nikeweke"
				return this.RenderTemplate("App/ws.html")
		}




		/*
		|-----------------------------------------
		|  RoomSocket
		|
		|---------------------------------------------
		*/
		func (this IndexController) RoomSocket(user string, ws *websocket.Conn) revel.Result {

						// Join the room.
						subscription := chatroom.Subscribe()
						defer subscription.Cancel()

						chatroom.Join(user)
						defer chatroom.Leave(user)

						// Send down the archive.
						for _, event := range subscription.Archive {
							if websocket.JSON.Send(ws, &event) != nil {
								// They disconnected
								return nil
							}
						}

						// In order to select between websocket messages and subscription events, we
						// need to stuff websocket events into a channel.
						newMessages := make(chan string)
						go func() {
							var msg string
							for {
								err := websocket.Message.Receive(ws, &msg)
								if err != nil {
									close(newMessages)
									return
								}
								newMessages <- msg
							}
						}()

						// Now listen for new events from either the websocket or the chatroom.
						for {
							select {
							case event := <-subscription.New:
								if websocket.JSON.Send(ws, &event) != nil {
									// They disconnected.
									return nil
								}
							case msg, ok := <-newMessages:
								// If the channel is closed, they disconnected.
								if !ok {
									return nil
								}

								// Otherwise, say something.
								chatroom.Say(user, msg)
							}
						}
						return nil
		}
