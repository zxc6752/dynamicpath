package lb_client

// import (
// 	// "crypto/tls"
// 	"dynamicpath/lib/loadbalancer_api"
// 	"dynamicpath/lib/up_topology"
// 	"dynamicpath/src/load_balancer/lb_context"
// 	"dynamicpath/src/load_balancer/logger"
// 	// "golang.org/x/net/http2"
// 	// "net/http"
// )

// func SendInitPathRequest() {
// 	// client := &http.Client{}
// 	// client.Transport = &http2.Transport{
// 	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	// }
// 	self := lb_context.LB_Self()
// 	// self.SmfClient = client
// 	request := loadbalancer_api.InitPathRequest{}
// 	for _, path := range self.Topology.PathListAll {
// 		pathInfo := loadbalancer_api.PathInfo{
// 			PathId: path.Id,
// 		}
// 		for _, edge := range path.Edges {
// 			edgeInfo := loadbalancer_api.EdgeInfo{
// 				StartIp:     edge.StartIp,
// 				EndIp:       edge.EndIp,
// 				StartNodeId: edge.Start.NodeId,
// 				EndNodeId:   edge.End.NodeId,
// 			}
// 			pathInfo.EdgeInfos = append(pathInfo.EdgeInfos, edgeInfo)
// 		}
// 		request.PathInfos = append(request.PathInfos, pathInfo)
// 	}
// 	// err := loadbalancer_api.InitPathListRequest(client, self.SmfUri, request)
// 	// if err != nil {
// 	// 	logger.SmfClientLog.Error(err.Error())
// 	// }
// }

// func SendWorstPathList(pathList []up_topology.Path) {
// 	request := loadbalancer_api.WorstPathList{}
// 	for _, path := range pathList {
// 		request.PathList = append(request.PathList, path.Id)
// 	}
// 	err := loadbalancer_api.UpdateWorstPathList(lb_context.LB_Self().SmfClient, lb_context.LB_Self().SmfUri, request)
// 	if err != nil {
// 		logger.SmfClientLog.Error(err.Error())
// 	}
// }

// func SendBestPath(supi string, request loadbalancer_api.SessionInfo) {
// 	err := loadbalancer_api.SendPduSessionPathRequest(lb_context.LB_Self().SmfClient, lb_context.LB_Self().SmfUri, supi, request)
// 	if err != nil {
// 		logger.SmfClientLog.Error(err.Error())
// 	}
// }
