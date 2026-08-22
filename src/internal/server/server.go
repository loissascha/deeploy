package server

import (
	"context"
	"encoding/json"
	"fmt"
	"local/deeploy/internal/communication"
	"log"
	"log/slog"
	"net/http"
	"slices"

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

func (s *Server) AddAgentIfNotExist(ctx context.Context, id int, conn *websocket.Conn) (*ServerAgent, error) {
	for _, a := range s.agents {
		if a.id == id {
			return nil, fmt.Errorf("agent with this 'id' is already registered")
		}
	}

	a := NewServerAgent(id, conn)
	err := a.SendWelcomeMessage(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("added new agent", "id", id)
	s.agents = append(s.agents, a)
	return a, nil
}

func (s *Server) RemoveAgent(a *ServerAgent) {
	slog.Info("removing agent", "id", a.id)
	s.agents = slices.DeleteFunc(s.agents, func(aa *ServerAgent) bool {
		if aa == a {
			if aa.conn != nil {
				aa.conn.CloseNow()
			}
			return true
		}
		return false
	})
}

func (s *Server) RunWS() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println("conn err", err)
			return
		}

		ctx := r.Context()

		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Println("read:", err)
			conn.CloseNow()
			return
		}

		var msg communication.WSMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println("unmarshal:", err)
			conn.CloseNow()
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
			a, err := s.AddAgentIfNotExist(ctx, msgInitialize.ID, conn)
			if err != nil {
				conn.Write(ctx, websocket.MessageText, []byte(err.Error())) // TODO: write error message or something instead?
				conn.CloseNow()
				return
			}

			a.RunReader(ctx)
			s.RemoveAgent(a)
		default:
			conn.CloseNow()
		}
	})

	log.Println("listening on :42066")
	err := http.ListenAndServe(":42066", nil)
	if err != nil {
		panic(err)
	}
}
