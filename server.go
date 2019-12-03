package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

type socket struct {
	conn        *websocket.Conn
	messagePipe chan []byte
}

const (
	writeWait  = 1 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WSServe(ws *socket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socketConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("WS failed: %s", err)))
			return
		}
		ws.conn = socketConn
		ws.messagePipe = make(chan []byte)

		go writePump(ws)
	}
}

func writePump(ws *socket) {
	newline := []byte{'\n'}
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		close(ws.messagePipe)
		fmt.Println("[WS] consumer stopped")
	}()
	for {
		select {
		case message, ok := <-ws.messagePipe:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ws.messagePipe)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-ws.messagePipe)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := DBConn()
	poe(err)

	data, err := GetAllItems(db)
	poe(err)

	err = db.Close()
	poe(err)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}

}

func Server() {
	router := mux.NewRouter()

	var sessWS *socket

	router.HandleFunc("/ws", WSServe(sessWS))
	router.HandleFunc("/", getAllHandler)

	// Run Server
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
		},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(router)
	bindAddr := fmt.Sprintf(":%d", 8081)

	log.Fatal(http.ListenAndServe(bindAddr, handler))
}
