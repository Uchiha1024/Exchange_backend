// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.1

package handler

import (
	"net/http"
	"github.com/zeromicro/go-zero/rest"
)

type Routes struct {
	server *rest.Server
	middleware []rest.Middleware
}

func NewRoutes(server *rest.Server) *Routes {
	return &Routes{
		server: server,
	}
}

func (r *Routes) Get(path string, handlerFunc http.HandlerFunc){


	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middleware,
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: handlerFunc,
			},
		),

	)

}


func (r *Routes) Post(path string, handlerFunc http.HandlerFunc){

	r.server.AddRoutes(
		rest.WithMiddlewares(
			r.middleware,
			rest.Route{
				Method:  http.MethodPost,
				Path:    path,
				Handler: handlerFunc,
			},
		),

	)

}


func (r *Routes) Group() *Routes{
	return &Routes{
		server: r.server,
	}
}

func (r *Routes) Use(middle ...rest.Middleware) {
	r.middleware = append(r.middleware, middle...)
}