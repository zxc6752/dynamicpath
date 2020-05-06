package main

import (
	"crypto/tls"
	"dynamicpath/lib/loadbalancer_api"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"time"
)

var host string
var times, period int
var client *http.Client
var pathList *loadbalancer_api.PathListAll

func Start(host string) {
	client = &http.Client{}
	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for {
		pathList, err := loadbalancer_api.GetPathListAll(client, host)
		if err == nil && pathList != nil {
			// for _, path := range pathList.PathInfos {
			// 	fmt.Printf("pathId: %d\n%s", path.PathId, path.EdgeInfos[0].StartNodeId)
			// 	for _, edge := range path.EdgeInfos {
			// 		fmt.Printf("\t%s", edge.EndNodeId)
			// 	}
			// 	fmt.Println()
			// }
			fmt.Println("[SMF] Start Load Balancer Success!!!")
			break
		}
		// fmt.Println("[SMF] Start Load Balancer Failed!!!")
	}

}

func main() {
	// runtime.GOMAXPROCS(10)
	flag.StringVar(&host, "h", "127.0.0.1:9487", "set `host` value")
	flag.IntVar(&times, "t", 5, "set `times` value")
	flag.IntVar(&period, "p", 3, "set `period` value")
	flag.Parse()

	Start(host)
	var sessionInfo []loadbalancer_api.SessionInfo
	for i := 0; i < times; i++ {
		req := loadbalancer_api.SessionInfo{
			PduSessionId: i + 1,
		}
		rsp, err := loadbalancer_api.SendPduSessionPathRequest(client, host, "imsi-123456789", req)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("%v\n", *rsp)
		sessionInfo = append(sessionInfo, *rsp)
		time.Sleep(time.Duration(period) * time.Second)
	}
	for i := 0; i < times; i++ {
		err := loadbalancer_api.SendPduSessionPathDeleteRequest(client, host, "imsi-123456789", sessionInfo[i])
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(time.Duration(period) * time.Second)
	}

	// router := http_server.NewRouter()
	// fmt.Printf("Server Listen on: %s\n", host)

	// mointorPemPath := path_relative.GetAbsPath("../../support/TLS/monitor.pem")
	// mointorKeyPath := path_relative.GetAbsPath("../../support/TLS/monitor.key")

	// server := &http.Server{
	// 	Addr:    host,
	// 	Handler: router,
	// }

	// err := server.ListenAndServeTLS(mointorPemPath, mointorKeyPath)
	// if err != nil {
	// 	log.Panic(err)
	// }
}
