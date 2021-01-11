package messages


import (
    "log"
	"encoding/json"
)



type ChatMessage struct {
    UserName string
    Message string
}


func JsonChatMessage(user_name string, message string) []byte {
    chat_message := ChatMessage{user_name, message}
    chat_message_json, err := json.Marshal(chat_message)
    if err != nil {
		log.Println("chat message encode error:", err, "message:", message, "user name:", user_name)
		return nil
	}
    return chat_message_json
}
