package request

import (
	model "bigagent/model/machine"
	"bigagent/util/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type PostStand struct {
	h    string
	c    *http.Client
	d    *model.StandData
	resp *http.Response
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewPostStand(host string) *PostStand {
	return &PostStand{h: host, c: &http.Client{}, d: model.NewStandData()}
}

func (p *PostStand) Do() (interface{}, error) {
	body, err := json.Marshal(p.d)
	if err != nil {
		logger.DefaultLogger.Error("Error marshaling JSON data: %v", err)
		return nil, err
	}

	resp, err := http.Post(p.h, "application/json", bytes.NewReader(body))
	if err != nil {
		logger.DefaultLogger.Error("Error making POST request to %s: %v", p.h, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.DefaultLogger.Error("Error closing response body: %v", err)
		}
	}(resp.Body)

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.DefaultLogger.Error("Error reading response body: %v", err)
		return nil, err
	}

	var response Response
	err = json.Unmarshal(bs, &response)
	if err != nil {
		logger.DefaultLogger.Error("Error unmarshaling response: %v", err)
		return nil, err
	}

	return response, nil
}
