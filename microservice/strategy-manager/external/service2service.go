package external

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/108356037/v1/strategy-manager/oauth"
	log "github.com/sirupsen/logrus"
)

// given the cpuOffset, memOffset, check if remain request resource is sufficient
func RequestValidation(user, cpu, mem string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var bearer = "Bearer " + oauth.BearerToken

	endpoint := fmt.Sprintf("http://%s/resources/%s/validate/requests", os.Getenv("USER_RESOURCE_SERVER_HOST"), user)

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Warn(err.Error())
		return false
	}

	param := req.URL.Query()
	param.Add("cpu_req", cpu)
	param.Add("mem_req", mem)
	req.URL.RawQuery = param.Encode()

	req.Header.Add("Authorization", bearer)

	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
		return false
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}

}

// given the cpuOffset, memOffset, check if remain limit resource is sufficient
func LimitValidation(user, cpu, mem string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var bearer = "Bearer " + oauth.BearerToken

	endpoint := fmt.Sprintf("http://%s/resources/%s/validate/limits", os.Getenv("USER_RESOURCE_SERVER_HOST"), user)

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Warn(err.Error())
		return false
	}

	param := req.URL.Query()
	param.Add("cpu_req", cpu)
	param.Add("mem_req", mem)
	req.URL.RawQuery = param.Encode()

	req.Header.Add("Authorization", bearer)

	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
		return false
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}

}
