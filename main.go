package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

    g.GET("/ping",func(ctx *gin.Context) {
		ctx.JSON(200,"pong")
	})

    g.GET("/dial",handle)

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	g.Run(fmt.Sprintf(":%s",port))
}

type Config struct{
	Route []*struct{
		Host string
		Path string
		Payload string
		Headers []string
		Method string
		Code int
		Body string
		Timeout int
	}
}

func handle(context *gin.Context) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "/data/config.toml"
	}
	var config Config
    if _, err := toml.DecodeFile(path,&config); err != nil {
		fmt.Printf("decode file %s failed: %+v\n",path,err)
		context.JSON(500,err.Error())
		return
	}


	for _, route := range config.Route {
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
		if err != nil {
			context.JSON(500,err.Error())
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			context.JSON(500,err.Error())
			return
		}
		if route.Code != 0 && route.Code != resp.StatusCode {
			context.JSON(500,fmt.Sprintf("%s: real code %d not equal expect code %d",dest,resp.StatusCode,route.Code))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			context.JSON(500,fmt.Sprintf("%s: read real response body failed: %+v",dest,err.Error()))
			return
		}
		if route.Body != "" && route.Body != string(body) {
			context.JSON(500,fmt.Sprintf("%s: real body %s not equal expect body %s",dest,string(body),route.Body))
			return
		}
	}

	fmt.Printf("%+v",config)
	context.JSON(200,config)
}