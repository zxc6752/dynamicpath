package lb_handler

import (
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/lb_util"
	"dynamicpath/src/load_balancer/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var HandlerLog *logrus.Entry

func init() {
	// init Pool
	HandlerLog = logger.HandlerLog
}

func Handle() {
	for {
		select {
		case msg, ok := <-lb_context.LBChannel:
			// Wait for Data Update
			if ok {
				lb_context.Wg.Wait()
				lb_context.Mtx.Lock()
				switch msg.Event {
				case lb_context.EventPduSessionAdd:
					val := msg.Value.(lb_context.PduSessionRequest)
					sessionKey := fmt.Sprintf("%s-%d", val.Supi, val.SessionInfo.PduSessionId)

					topology := lb_context.LB_Self().Topology
					session := lb_context.PduSessionContext{
						Supi:        val.Supi,
						SessionInfo: val.SessionInfo,
					}
					bestPathId := topology.RoundRobinCnt

					if lb_context.LB_Self().LoadBalancerType == lb_util.LoadBalancerType_DP {
						allOverload := true
						for i := 0; i < lb_util.DP_ShortestPathThreshlod; i++ {
							if lb_util.AvailableBandwidth[i] < 0 {
								continue
							}
							lb_util.AvailableBandwidth[i] -= topology.Granularity
							allOverload = false
							bestPathId = topology.PathListAll[i].Id
						}
						if allOverload {
							lb_util.DP_ShortestPathThreshlod = -2
							bestPathId = topology.PathListAll[0].Id
						}
					} else {

						bestList := topology.PathListBest
						if bestList == nil {
							// Round-Robin if all path overflow
							topology.RoundRobinCnt = (topology.RoundRobinCnt + 1) % len(topology.PathInfos)
						} else {
							// Decide Best Path And Send to SMF
							bestPathId = bestList[0].Id
							lb_util.ModifyPathRemainRate(bestList[0], lb_util.MinusRemainRate, topology.Granularity)
						}
					}
					logger.UtilLog.Debugf("Supi-SessionId: %s, PathId: %d Load: %.2f", sessionKey, bestPathId, *topology.PathInfos[bestPathId].Path.RemainRate)

					session.SessionInfo.PathID = bestPathId
					topology.PathInfos[bestPathId].PduSessionInfo[sessionKey] = &session
					lb_context.SendHttpResponseMessage(msg.HttpChannel, nil, http.StatusCreated, session.SessionInfo)

				case lb_context.EventPduSessionDel:
					val := msg.Value.(lb_context.PduSessionRequest)
					sessionKey := fmt.Sprintf("%s-%d", val.Supi, val.SessionInfo.PduSessionId)

					topology := lb_context.LB_Self().Topology
					pathId := val.SessionInfo.PathID
					pathInfo := topology.PathInfos[pathId]
					if pathInfo != nil {
						lb_context.SendHttpResponseMessage(msg.HttpChannel, nil, http.StatusNoContent, nil)
						lb_util.ModifyPathRemainRate(pathInfo.Path, lb_util.AddRemainRate, topology.Granularity)
						delete(topology.PathInfos[pathId].PduSessionInfo, sessionKey)
					} else {
						lb_context.SendHttpResponseMessage(msg.HttpChannel, nil, http.StatusNotFound, nil)
						logger.HandlerLog.Warnf("pathId: %d Not Found", val.SessionInfo.PathID)
					}
				}
				lb_context.Mtx.Unlock()

			} else {
				HandlerLog.Errorln("Channel closed!")
			}
		case <-time.After(time.Second * 1):

		}
	}
}
