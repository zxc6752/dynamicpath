package http_server

import (
	"dynamicpath/lib/loadbalancer_api"
	"dynamicpath/src/load_balancer/lb_context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllPath(c *gin.Context) {
	self := lb_context.LB_Self()
	// self.SmfClient = client
	resp := loadbalancer_api.PathListAll{}
	for _, path := range self.Topology.PathListAll {
		pathInfo := loadbalancer_api.PathInfo{
			PathId: path.Id,
		}
		for _, edge := range path.Edges {
			edgeInfo := loadbalancer_api.EdgeInfo{
				StartIp:     edge.StartIp,
				EndIp:       edge.EndIp,
				StartNodeId: edge.Start.NodeId,
				EndNodeId:   edge.End.NodeId,
			}
			pathInfo.EdgeInfos = append(pathInfo.EdgeInfos, edgeInfo)
		}
		resp.PathInfos = append(resp.PathInfos, pathInfo)
	}
	c.JSON(http.StatusOK, resp)
}

func PostPduSessionPath(c *gin.Context) {
	var sessionInfo loadbalancer_api.SessionInfo
	if err := c.ShouldBindJSON(&sessionInfo); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	val := lb_context.PduSessionRequest{
		Supi:        c.Params.ByName("supi"),
		SessionInfo: sessionInfo,
	}
	channelMsg := lb_context.NewHttpChannelMessage(lb_context.EventPduSessionAdd, val)
	lb_context.SendMessage(channelMsg)
	recvMsg := <-channelMsg.HttpChannel

	c.JSON(recvMsg.HTTPResponse.Status, recvMsg.HTTPResponse.Body)
}

func DelPduSessionPath(c *gin.Context) {
	var sessionInfo loadbalancer_api.SessionInfo
	if err := c.ShouldBindJSON(&sessionInfo); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		fmt.Println("[ERROR]", err.Error())
		return
	}
	val := lb_context.PduSessionRequest{
		Supi:        c.Params.ByName("supi"),
		SessionInfo: sessionInfo,
	}
	channelMsg := lb_context.NewHttpChannelMessage(lb_context.EventPduSessionDel, val)
	lb_context.SendMessage(channelMsg)
	recvMsg := <-channelMsg.HttpChannel

	c.JSON(recvMsg.HTTPResponse.Status, recvMsg.HTTPResponse.Body)
}
