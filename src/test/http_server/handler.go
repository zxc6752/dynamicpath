package http_server

import (
	"dynamicpath/lib/loadbalancer_api"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostAllPath(c *gin.Context) {
	var request loadbalancer_api.InitPathRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	spew.Dump(request)
	c.JSON(http.StatusNoContent, gin.H{})
}

func PostWorstPath(c *gin.Context) {
	var request loadbalancer_api.WorstPathList
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	spew.Dump(request)
	c.JSON(http.StatusNoContent, gin.H{})
}

func PostPduSessionPath(c *gin.Context) {
	var request loadbalancer_api.SessionInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	spew.Dump(request)
	c.JSON(http.StatusNoContent, gin.H{})
}
