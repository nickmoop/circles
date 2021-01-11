package handlers

import (
    "flag"
	"log"
	"net/http"
)


func ApplyHandlers(hostname *string) {
    flag.Parse()
	log.SetFlags(0)

    // files providers
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("../js"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../css"))))

    // http handlers
	http.HandleFunc("/", home)
	http.HandleFunc("/createSession", createSession)

	// web socket handlers
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/chatJoin", chatWSSessionJoin)
	http.HandleFunc("/battleJoin", battleWSSessionJoin)
	http.HandleFunc("/battle", battleHandler)

	log.Println("Server start listening:", *hostname)
	log.Fatal(http.ListenAndServe(*hostname, nil))
}
