package lb_client

import (
	"crypto/tls"
	"dynamicpath/lib/MonitorInfo"
	"dynamicpath/lib/monitor_api"
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/logger"
	"golang.org/x/net/http2"
	"net/http"
)

func StartMonitor(upf *lb_context.UpfContext) {
	client := &http.Client{}
	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	upf.HttpClient = client
	request := MonitorInfo.MonitorStartRequest{}
	for _, edge := range upf.TopologyNode.Edges {
		info := MonitorInfo.MonitorIPInfo{
			IP: edge.EndIp,
		}
		request.LinkIPInfo = append(request.LinkIPInfo, info)
	}
	if upf.RemoteIpToEdge != nil {
		for i, info := range request.LinkIPInfo {
			if _, exist := upf.RemoteIpToEdge[info.IP]; exist {
				request.LinkIPInfo[i].Update = true
			}
		}
	}
	respInfo, err := monitor_api.StartMonitorRequest(client, upf.HttpServerUri, request)
	if err != nil {
		logger.MonitorClientLog.Error(err.Error())
		return
	}
	upf.RemoteIpMaps = respInfo.RemoteIpMaps
	for ethName, bandwidth := range respInfo.BandwidthMaps {
		upf.InterfaceInfos[ethName] = lb_context.InterfaceInfo{
			TotalBandwidth: bandwidth,
			EthName:        ethName,
		}
		rx, tx := bandwidth.Rx, bandwidth.Tx
		upf.TopologyNode.RemainRates[lb_context.GetIfaceRxName(ethName)] = &rx
		upf.TopologyNode.RemainRates[lb_context.GetIfaceTxName(ethName)] = &tx
	}
}

func GetMonitorData(upf *lb_context.UpfContext, delta float64, updateCost bool) {
	respInfo, err := monitor_api.GetMonitorData(upf.HttpClient, upf.HttpServerUri)
	if err != nil {
		logger.MonitorClientLog.Error(err.Error())
		return
	}
	if updateCost {
		for remoteIp, connectionInfo := range respInfo.ConnectionInfos {
			edge := upf.RemoteIpToEdge[remoteIp]
			logger.MonitorClientLog.Tracef("%s -> %s, delay: %.2fms, loss: %.2f%%, Cost: %.2f", edge.StartIp, edge.EndIp, connectionInfo.DelayTime, connectionInfo.PacketLoss, getCostValue(delta, connectionInfo.DelayTime, connectionInfo.PacketLoss))
			edge.UpdateCost(getCostValue(delta, connectionInfo.DelayTime, connectionInfo.PacketLoss))
		}
	}
	for ethName, bandwidth := range respInfo.PacketRates {
		ifaceInfo := upf.InterfaceInfos[ethName]
		logger.MonitorClientLog.Tracef("NodeId: %s\tEthName: %s\tTotalRx: %.2fmbps\tCurrRx: %.2fmbps\tTotalTx: %.2fmbps\tCurrTx: %.2fmbps",
			upf.TopologyNode.NodeId, ethName, ifaceInfo.TotalBandwidth.Rx, bandwidth.Rx, ifaceInfo.TotalBandwidth.Tx, bandwidth.Tx)
		*upf.TopologyNode.RemainRates[lb_context.GetIfaceRxName(ethName)] = ifaceInfo.TotalBandwidth.Rx - bandwidth.Rx
		*upf.TopologyNode.RemainRates[lb_context.GetIfaceTxName(ethName)] = ifaceInfo.TotalBandwidth.Tx - bandwidth.Tx
	}

}

func getCostValue(delta, delay, packetlossRate float64) float64 {
	if packetlossRate == 1.0 {
		return 100000.0 // if all packet loss, then Cost = 10 seconds
	}
	return (delta + (1-delta)*(1+packetlossRate)/(1-packetlossRate)) * delay
}
