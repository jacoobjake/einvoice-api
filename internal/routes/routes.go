package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stephenafamo/bob"
)

type RouteRegistrar struct {
	r  *gin.Engine
	db *bob.DB
}

func NewRouteRegistrar(r *gin.Engine, db *bob.DB) *RouteRegistrar {
	return &RouteRegistrar{
		r:  r,
		db: db,
	}
}

func (registrar RouteRegistrar) RegisterRoutes() {
	registrar.r.GET("/test")
}
