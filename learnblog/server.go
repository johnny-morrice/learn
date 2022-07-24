package main

import "github.com/gin-gonic/gin"

type Router interface {
	AddRoutesToGroup(group *gin.RouterGroup)
}

type RouterPackage struct {
	BasePath string
	Router   Router
}

type Server struct {
	RouterPackages []RouterPackage
	Addr           string
}

func (server *Server) Run() error {
	r := gin.Default()
	for _, pack := range server.RouterPackages {
		group := r.Group(pack.BasePath)
		pack.Router.AddRoutesToGroup(group)
	}
	return r.Run(server.Addr)
}
