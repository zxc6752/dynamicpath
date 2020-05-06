package loadbalancer_api

import (
	"bytes"
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

func GetPathListAll(client *http.Client, host string) (response *PathListAll, err error) {
	resp, err := client.Get(fmt.Sprintf("https://%s/loadbalancer/init", host))
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

func SendPduSessionPathRequest(client *http.Client, host string, supi string, request SessionInfo) (response *SessionInfo, err error) {
	body, err := prepareRequest(request)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post(fmt.Sprintf("https://%s/loadbalancer/user/%s", host, supi), "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf(resp.Status)
	}
	rspbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rspbody, &response)
	return
}

func SendPduSessionPathDeleteRequest(client *http.Client, host string, supi string, request SessionInfo) error {
	body, err := prepareRequest(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/loadbalancer/user/%s", host, supi), body)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf(resp.Status)
	}
	return nil
}

// func UpdateWorstPathList(client *http.Client, host string, request WorstPathList) error {
// 	body, err := prepareRequest(request)
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := client.Post(fmt.Sprintf("https://%s/loadbalancer/path", host), "application/json", body)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusNoContent {
// 		return fmt.Errorf(resp.Status)
// 	}
// 	return nil
// }
