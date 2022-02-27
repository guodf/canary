package controller

import "github.com/gin-gonic/gin"

type IHttpController interface {
	Register(*gin.Engine)
}
