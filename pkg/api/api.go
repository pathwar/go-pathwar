package api

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

var (
	APIUrl = "https://api.pathwar.net"
)

type APIPathwar struct {
	token string
	debug bool
}

func NewAPIPathwar(token, debug string) *APIPathwar {
	return &APIPathwar{
		token: token,
		debug: debug != "",
	}
}

func printErrors(errs []error) error {
	for _, err := range errs {
		logrus.Error(err)
	}
	return errors.New("Error(s) has occured")
}

func httpHandleError(goodStatusCode []int, statusCode int, body []byte) error {
	good := false
	for _, code := range goodStatusCode {
		if code == statusCode {
			good = true
		}
	}
	if !good {
		return errors.New(string(body))
	}
	return nil
}

func (p *APIPathwar) GetRequest(url string) ([]byte, error) {
	request := gorequest.New().Get(strings.Join([]string{APIUrl, url}, "/"))
	request = request.SetBasicAuth(p.token, "")
	if p.debug {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{200}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (p *APIPathwar) DeleteRequest(url, etag string) ([]byte, error) {
	request := gorequest.New().Delete(strings.Join([]string{APIUrl, url}, "/"))
	request = request.SetBasicAuth(p.token, "")
	request = request.Set("If-Match", etag)
	if p.debug {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{200}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (p *APIPathwar) PatchRequest(url, etag string, data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request := gorequest.New().Patch(strings.Join([]string{APIUrl, url}, "/"))
	request = request.SetBasicAuth(p.token, "")
	request = request.Set("If-Match", etag)
	request = request.Send(string(payload))
	if p.debug {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{200}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (p *APIPathwar) PostRequest(url string, data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request := gorequest.New().Post(strings.Join([]string{APIUrl, url}, "/"))
	request = request.SetBasicAuth(p.token, "")
	request = request.Send(string(payload))
	if p.debug {
		request = request.SetDebug(true)
	}
	resp, body, errs := request.EndBytes()

	if len(errs) != 0 {
		return nil, printErrors(errs)
	}
	if err := httpHandleError([]int{201}, resp.StatusCode, body); err != nil {
		return nil, err
	}
	return body, nil
}
