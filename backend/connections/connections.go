package connections


import (
	"log"
	"net/http"
    "sync"

    "github.com/gorilla/websocket"
)


var upgrader = getUpgrader()


type WSConnection struct {
    SessionId string
    UserName string
    Connection *websocket.Conn
    Mutex sync.Mutex
}


func NewWSConnection(writer http.ResponseWriter, request *http.Request) (*WSConnection){
    connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("connection creation error:", err)
		return nil
	}

    ws_connection := new(WSConnection)
	ws_connection.Connection = connection

    return ws_connection
}


func (connection *WSConnection) ReadMessage() ([]byte, error) {
    _, message, err := connection.Connection.ReadMessage()
    if err != nil {
        log.Println("read error:", err)
    }

    return message, err
}


func (connection *WSConnection) ReadMessageAsString() (string, error) {
    message, err := connection.ReadMessage()

    return string(message), err
}


func (connection *WSConnection) WriteMessage(message []byte) {
    connection.Mutex.Lock()
    defer connection.Mutex.Unlock()
    err := connection.Connection.WriteMessage(1, message)
    if err != nil {
        log.Println("write error:", err)
    }
}


func (connection *WSConnection) Close() {
    log.Println("connection closed for session:", connection.SessionId, "and user:", connection.UserName)
    connection.Connection.Close()
    connection.Connection = nil
}


func getUpgrader() (websocket.Upgrader) {
    var upgrader = websocket.Upgrader{}
    upgrader.CheckOrigin = func(request *http.Request) bool { return true }

    return upgrader
}
