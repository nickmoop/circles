package handlers

import (
	"log"
	"encoding/json"
	"html/template"
	"net/http"

	"backend/sessions"
)



func createSession(writer http.ResponseWriter, request *http.Request) {
    session_json := sessions.NewCoreSessionJson(request)

    log.Println("create new session:", session_json.SessionId, "for user:", session_json.UserName)

    if err := sessions.TMPAddToSessionsPool(session_json); err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(*err))
    } else {
        json.NewEncoder(writer).Encode(session_json)
    }

}


func home(writer http.ResponseWriter, request *http.Request) {

    var TMPhomeTemplate, err = template.ParseFiles("../html/main.html")

    if err != nil {
        log.Println("error while read template:", err)
    }

    TMPhomeTemplate.Execute(writer, nil)
}
