package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		fmt.Printf("Panic occurred: %s \nstacktrace : %s", recovered, string(debug.Stack()))
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	return r
}
