package http_server_test

import (
	"crypto/tls"
	"dynamicpath/lib/MonitorInfo"
	"dynamicpath/lib/monitor_api"
	"dynamicpath/lib/path_util"
	"dynamicpath/src/monitor/http_server"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func Init() {
	router := http_server.NewRouter()

	monitorLogPath := path_util.ModulePath("dynamicpath/support/TLS/monitor.log")
	mointorPemPath := path_util.ModulePath("dynamicpath/support/TLS/monitor.pem")
	mointorKeyPath := path_util.ModulePath("dynamicpath/support/TLS/monitor.key")

	keylogFile, err := os.OpenFile(monitorLogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Panic(err)
	}
	server := &http.Server{
		Addr: ":8888",
		TLSConfig: &tls.Config{
			KeyLogWriter: keylogFile,
		},
		Handler: router,
	}

	go server.ListenAndServeTLS(mointorPemPath, mointorKeyPath)

	time.Sleep(100 * time.Millisecond)

}

func TestAll(t *testing.T) {
	client := &http.Client{}
	client.Transport = &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	{

		request := MonitorInfo.MonitorStartRequest{
			LinkIPInfo: []MonitorInfo.MonitorIPInfo{
				{
					IP:     "8.8.8.8",
					Update: true,
				},
				{
					IP:     "8.8.4.4",
					Update: true,
				},
				{
					IP:     "192.188.2.4",
					Update: true,
				},
			},
		}
		{
			respInfo, err := monitor_api.StartMonitorRequest(client, "localhost:8888", request)
			assert.True(t, err == nil)
			spew.Dump(respInfo)
		}
		{
			_, err := monitor_api.StartMonitorRequest(client, "localhost:8888", request)
			assert.True(t, err != nil)
			fmt.Println(err.Error())
		}

	}
	time.Sleep(15 * time.Second)
	{
		respInfo, err := monitor_api.GetMonitorData(client, "localhost:8888")
		assert.True(t, err == nil)
		spew.Dump(respInfo)
	}

}
