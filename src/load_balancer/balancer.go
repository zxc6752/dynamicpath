package main

import (
	"crypto/tls"
	"dynamicpath/lib/path_util"
	"dynamicpath/src/load_balancer/factory"
	"dynamicpath/src/load_balancer/http_server"
	"dynamicpath/src/load_balancer/lb_client"
	"dynamicpath/src/load_balancer/lb_context"
	"dynamicpath/src/load_balancer/lb_handler"
	"dynamicpath/src/load_balancer/lb_util"
	"dynamicpath/src/load_balancer/logger"
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
)

var config string

func Initailize() {
	flag.StringVar(&config, "lbcfg", path_util.ModulePath("dynamicpath/config/lbcfg.conf"), "load balancer config file")
	flag.Parse()

	factory.InitConfigFactory(config)
	config := factory.LBConfig
	if config.Logger.DebugLevel != "" {
		level, err := logrus.ParseLevel(config.Logger.DebugLevel)
		if err == nil {
			logger.SetLogLevel(level)
		}
	}
	logger.SetReportCaller(config.Logger.ReportCaller)

}

func main() {
	Initailize()
	self := lb_context.LB_Self()
	lb_util.InitTopology()
	wg := sync.WaitGroup{}
	wg.Add(len(self.UPFInfos))
	for _, upf := range self.UPFInfos {
		go func(upf *lb_context.UpfContext) {
			lb_client.StartMonitor(upf)
			wg.Done()
		}(upf)
	}
	wg.Wait()
	self.LoadBalancerType = factory.LBConfig.LoadBalancerType
	lb_util.InitAllPath()
	// Send PathList to Smf
	// lb_client.SendInitPathRequest()
	// topology := self.Topology
	// topology.PathThresh = lb_util.FindPathThreshold(topology.PathListAll)
	// lb_util.UpdatePathRemainRate(topology.PathListAll[:topology.PathThresh+1])
	lb_util.RefreshAll(factory.LBConfig.Period)

	go lb_handler.Handle()

	router := http_server.NewRouter()
	logger.InitLog.Infof("Server Listen on: %s\n", self.Host)

	monitorLogPath := path_util.ModulePath("dynamicpath/support/TLS/loadbalancer.log")
	mointorPemPath := path_util.ModulePath("dynamicpath/support/TLS/loadbalancer.pem")
	mointorKeyPath := path_util.ModulePath("dynamicpath/support/TLS/loadbalancer.key")

	keylogFile, err := os.OpenFile(monitorLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.InitLog.Error(err.Error())
	}
	server := &http.Server{
		Addr: self.Host,
		TLSConfig: &tls.Config{
			KeyLogWriter: keylogFile,
		},
		Handler: router,
	}

	err = server.ListenAndServeTLS(mointorPemPath, mointorKeyPath)
	if err != nil {
		logger.InitLog.Error(err.Error())
	}
}
