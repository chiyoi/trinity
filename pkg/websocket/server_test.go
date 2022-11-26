package websocket

import (
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	var wg sync.WaitGroup
	go func() {
		wsConfig := ServerConfig{}
		srv := http.Server{}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			ws, err := wsConfig.Hijack(w, r)
			if err != nil {
				log.Println("server: error:", err)
				return
			}
			conn := ws.(*webSocket).conn
			defer func() {
				err := ws.Close()
				if err != nil {
					log.Println(err)
				}
			}()
			log.Println("server:", ws.State())

			for c := 2; c > 0; c-- {
				time.Sleep(time.Second * 3)
				_, err = conn.Write([]byte("ping"))
				if err != nil {
					log.Println("server error:", err)
				}
				log.Println("server sent:", "ping")

				var b = make([]byte, 4)
				_, err = conn.Read(b)
				if err != nil {
					log.Println("server error:", err)
				}
				log.Println("server received:", string(b))
			}
			wg.Done()
		})
		log.Fatalln(srv.ListenAndServe())
	}()
	go func() {
		time.Sleep(time.Second * 3)
		ws, err := Dial("ws://localhost/")
		if err != nil {
			log.Fatalln(err)
		}
		defer func(ws WebSocket) {
			err := ws.Close()
			if err != nil {
				log.Println(err)
			}
		}(ws)
		log.Println("client:", ws.State())
		conn := ws.(*webSocket).conn

		for c := 2; c > 0; c-- {
			var b = make([]byte, 4)
			_, err = conn.Read(b)
			if err != nil {
				log.Println("client error:", err)
			}
			log.Println("client received:", string(b))

			_, err = conn.Write([]byte("pong"))
			if err != nil {
				log.Println("client error:", err)
			}
			log.Println("client sent:", "pong")
		}
		wg.Done()
	}()
	wg.Add(2)
	wg.Wait()
}
