package request

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type PostStand struct {
	h string
	c *http.Client
	p *bytes.Buffer
}

func NewPostStand(host string) *PostStand {
	return &PostStand{h: host, c: &http.Client{}, p: bytes.NewBufferString("")}
}

func (p *PostStand) Do() (interface{}, error) {
	resp, err := p.c.Post(p.h, "application/x-www-form-urlencoded", p.p)
	if err != nil {
		log.Printf("Error making request to %s with body %v: %v", p.h, p.p, err)
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v", p.h, err)
		return nil, err
	}
	return string(body), nil
}
