package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type CusDispatch struct {
	Config *Config
	Context *gin.Context
}

func (d *CusDispatch) Dispatch(){
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
			d.Context.JSON(500,err.Error())
			return 
		}
		if route.Code != 0 && route.Code != resp.StatusCode {
			d.Context.JSON(500,fmt.Sprintf("%s: real code %d not equal expect code %d",dest,resp.StatusCode,route.Code))
			return 
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			d.Context.JSON(500,fmt.Sprintf("%s: read real response body failed: %+v",dest,err.Error()))
			return 
		}
		if route.Body != "" && route.Body != string(body) {
			d.Context.JSON(500,fmt.Sprintf("%s: real body %s not equal expect body %s",dest,string(body),route.Body))
			return 
		}
	}
	d.Context.JSON(200,d.Config)
	return 
}