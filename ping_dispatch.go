package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type PingDispatch struct {
	Config *Config
	Context *gin.Context
}

func (d *PingDispatch) Dispatch() {
	for _, route := range d.Config.Route {
		client := http.DefaultClient
		if route.Timeout != 0 {
			client.Timeout = time.Duration(route.Timeout) * time.Second
		}else {
			client.Timeout = 10 * time.Second
		}
	    dest := fmt.Sprintf("%s%s", route.Host, route.Path)
		hd := make(http.Header,0)
		for _, header := range route.Headers {
			h := strings.Split(header, ": ")
			hd.Add(h[0],h[1])
		}
		req ,err := http.NewRequest(route.Method,dest,bytes.NewReader([]byte(route.Payload)))
		req.Header = hd
		if err != nil {
			d.Context.JSON(500,err.Error())
			return 
		}
		resp, err := client.Do(req)
		if err != nil {
			d.Context.JSON(500,"failed")
			return 
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d.Context.JSON(500,fmt.Sprintf("%s: read real response body failed: %+v",dest,err.Error()))
			return 
		}
		if route.Body != "" && route.Body != string(body) {
			d.Context.JSON(500,"failed")
			return 
		}
	}
	d.Context.JSON(200,"success")
	return 
}