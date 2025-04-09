package handler

import (
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

type Routers struct {
	server      *rest.Server
	middlewares []rest.Middleware
	prefix      string
}

func NewRouters(server *rest.Server, prefix string) *Routers {
	return &Routers{
		server: server,
		prefix: prefix,
	}
}

func (r *Routers) Get(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: handlerFunc,
			},
		),
		rest.WithPrefix(r.prefix),
	)
}

func (r *Routers) Post(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodPost,
				Path:    path,
				Handler: handlerFunc,
			},
		),
		rest.WithPrefix(r.prefix),
	)
}

func (r *Routers) GetNoPrefix(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: handlerFunc,
			},
		),
	)
}

func (r *Routers) PostNoPrefix(path string, handlerFunc http.HandlerFunc) {
	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middlewares,
			rest.Route{
				Method:  http.MethodPost,
				Path:    path,
				Handler: handlerFunc,
			},
		),
	)
}
func (r *Routers) Group() *Routers {
	return &Routers{
		server: r.server,
		prefix: r.prefix,
	}
}
