package ws

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strings"
)

const ROOM = "market"

type WebsocketServer struct {
	path   string
	server *socketio.Server
}

func (ws *WebsocketServer) Start() {
	ws.server.Serve()
}

func (ws *WebsocketServer) Stop() {
	ws.server.Close()
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func NewWebsocketServer(path string) *WebsocketServer {
	//解决跨域
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		logx.Info("connected:", s.ID())
		s.Join(ROOM)
		return nil
	})
	return &WebsocketServer{
		path:   path,
		server: server,
	}
}

// BroadcastToNamespace  "/" "/topic/market/thumb"
func (w *WebsocketServer) BroadcastToNamespace(path string, event string, data any) {
	go func() {
		w.server.BroadcastToRoom(path, ROOM, event, data)
	}()
}

func (ws *WebsocketServer) ServerHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//所有的请求都要通过这里
		path := r.URL.Path
		logx.Info("==========================", path)
		if strings.HasPrefix(path, ws.path) {
			//进行我们的处理
			ws.server.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
