package server

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Server struct {
	agents []*ServerAgent
}

type ServerAgent struct {
	id   int
	conn *websocket.Conn
}

func NewServer() *Server {
	return &Server{
		agents: []*ServerAgent{},
	}
}

func NewServerAgent(id int, conn *websocket.Conn) *ServerAgent {
	return &ServerAgent{
		id:   id,
		conn: conn,
	}
}

func (s *Server) RunWS() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println("conn err", err)
			return
		}
		defer conn.CloseNow()

		ctx := r.Context()

		for {
			msgType, data, err := conn.Read(ctx)
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("received: %s\n", data)

			err = conn.Write(ctx, msgType, data)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	})

	log.Println("listening on :42066")
	err := http.ListenAndServe(":42066", nil)
	if err != nil {
		panic(err)
	}
}
