package handlers

import (
	"log"
	"net/http"

    "github.com/gorilla/websocket"

    "backend/battle"
    "backend/connections"
    "backend/sessions"
)


var upgrader = getUpgrader()


func chatWSSessionJoin(writer http.ResponseWriter, request *http.Request) {
    connection := connections.NewWSConnection(writer, request)
    defer connection.Close()

    sessions.JoinChatConnectionToSession(connection)
    session_id := connection.SessionId
    user_name := connection.UserName

    sessions.TMPSendChatMessageToAllInSession(session_id, user_name, "has joined to chat.")

    log.Println("start listening chat messages for session:", connection.SessionId, "user name:", connection.UserName)
    for {
		message, err := connection.ReadMessageAsString()

		if err != nil {
		    log.Println("error on reading chat message:", err)
		    sessions.RemoveChatConnectionFromSession(connection)
		    break
		}

		sessions.TMPSendChatMessageToAllInSession(session_id, user_name, message)
	}

	sessions.TMPSendChatMessageToAllInSession(session_id, user_name, "has left chat.")
}


func battleWSSessionJoin(writer http.ResponseWriter, request *http.Request) {
    connection := connections.NewWSConnection(writer, request)
    defer connection.Close()

    sessions.JoinBattleConnectionToSession(connection)
    session_id := connection.SessionId
    user_name := connection.UserName

    sessions.TMPSendChatMessageToAllInSession(session_id, user_name, "has joined to battle.")

    log.Println("start listening battle messages for session:", connection.SessionId, "user name:", connection.UserName)
    for {
        message, err := connection.ReadMessageAsString()

		if err != nil {
		    log.Println("error on reading battle message:", err)
		    sessions.RemoveChatConnectionFromSession(connection)
		    break
		}

        log.Println("battle message:", message)
        response := battle.ExecuteBattleCommand(user_name, session_id, message)
        sessions.TMPSendBattleMessageToAllInSession(session_id, user_name, response)
	}
}


func battleHandler(writer http.ResponseWriter, request *http.Request) {
    connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		message_type, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(message_type, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}


func echo(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		message_type, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(message_type, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}


func getUpgrader() (websocket.Upgrader) {
    var upgrader = websocket.Upgrader{}
    upgrader.CheckOrigin = func(request *http.Request) bool { return true }

    return upgrader
}
