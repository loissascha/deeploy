package server

import (
	"context"
	"encoding/json"
	"fmt"
	"local/deeploy/internal/communication"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type Server struct {
	agents []*ServerAgent
}

func NewServer() *Server {
	return &Server{
		agents: []*ServerAgent{},
	}
}

func (s *Server) AddAgentIfNotExist(ctx context.Context, id int, conn *websocket.Conn) error {
	for _, a := range s.agents {
		if a.id == id {
			return fmt.Errorf("agent with this 'id' is already registered")
		}
	}

	a := NewServerAgent(id, conn)
	err := a.SendWelcomeMessage(ctx)
	if err != nil {
		return err
	}
	s.agents = append(s.agents, a)
	return nil
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
			_, data, err := conn.Read(ctx)
			if err != nil {
				log.Println("read:", err)
				return
			}

			var msg communication.WSMessage
			err = json.Unmarshal(data, &msg)
			if err != nil {
				log.Println("unmarshal:", err)
				return
			}

			switch msg.MsgType {
			case communication.MsgTypeInitiate:
				log.Println("initialize message received")
				var msgInitialize communication.InitiateMessage
				err = json.Unmarshal(msg.Data, &msgInitialize)
				if err != nil {
					log.Println("unmarshal 2:", err)
					return
				}
				err = s.AddAgentIfNotExist(ctx, msgInitialize.ID, conn)
				if err != nil {
					err = conn.Write(ctx, websocket.MessageText, []byte(err.Error())) // TODO: write error message or something instead?
					if err != nil {
						log.Println("write:", err)
						return
					}
					return
				}
				break
			}
		}
	})

	log.Println("listening on :42066")
	err := http.ListenAndServe(":42066", nil)
	if err != nil {
		panic(err)
	}
}
