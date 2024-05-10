package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("started")

	chat := NewChatroom()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("new websocket connection")
		screenName := r.URL.Query().Get("name")
		if screenName == "" {
			screenName = "anon"
		}
		NewClient(screenName, &chat, w, r)
		time.Sleep(200 * time.Millisecond)
		chat.messages <- fmt.Sprintf("=== %s joined the chat ===", screenName)
	})

	fmt.Println("listening on :8000")
	http.ListenAndServe(":8000", nil)
}
