package socket_server

import (
	"app_chat/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
)

type Server struct {
	Port            int
	Name            string
	ReadTimeOut     int
	WriteTimeOut    int
	MaxMessageSize  int
	ReadBufferSize  int
	WriteBufferSize int
	httpServer      *http.Server
	Handler         Handler
}

func (s *Server) Start() error {
	// TEST WS
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck/_check", handleHealthCheck).Methods("GET")
	r.PathPrefix("/").HandlerFunc(s.serveWebSocket)
	fmt.Printf("Server %v start at port :%v", s.Name, s.Port)
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "PUT", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%+v", s.Port),
		Handler: handlers.CORS(credentials, methods, origins)(r),
	}
	s.httpServer = srv
	return srv.ListenAndServe()
}
func (s *Server) StopServer() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Error stop http server: %v", err)
	} else {
		fmt.Printf("Http server %v stopped", s.Name)
	}
}

func handleHealthCheck(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode("OK")
}

func (s *Server) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  s.ReadBufferSize,
		WriteBufferSize: s.WriteBufferSize,
	}
	fmt.Printf("Serve connection %+v", r)
	conn, err := upgrade.Upgrade(w, r, nil)
	conn.SetReadLimit(int64(s.MaxMessageSize))
	if err != nil {
		fmt.Printf("Error while upgrade connection %v", err)
		return
	}
	utils.RunWithRecovery(func() {
		err := s.Handler.HandleSocketConnection(conn, r, s.ReadTimeOut, s.WriteTimeOut)
		if err != nil {
			return
		}
	})
}
