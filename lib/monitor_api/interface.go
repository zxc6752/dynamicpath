package monitor_api

import (
	"bytes"
	"dynamicpath/lib/MonitorInfo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func prepareRequest(body interface{}) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}
	err = json.NewEncoder(bodyBuf).Encode(body)
	return
}

func GetMonitorData(client *http.Client, host string) (response *MonitorInfo.UPFMonitorData, err error) {
	resp, err := client.Get(fmt.Sprintf("https://%s/monitor", host))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}
	rspbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rspbody, &response)
	return
}

func StartMonitorRequest(client *http.Client, host string, request MonitorInfo.MonitorStartRequest) (response *MonitorInfo.InterfacesInfo, err error) {
	body, err := prepareRequest(request)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(fmt.Sprintf("https://%s/monitor/start", host), "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}
	rspbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rspbody, &response)
	return
}
