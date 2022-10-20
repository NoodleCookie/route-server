package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"os"

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
		Code int
		Headers []string
		Body string
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
	    dest := fmt.Sprintf("%s%s", route.Host, route.Path)
		resp, err := client.Get(dest)
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
		for _, header := range route.Headers {
			h := strings.Split(header, ": ")
			key := h[0]
			value := h[1]
			v, ok := resp.Header[key]
			if !ok {
				context.JSON(500,fmt.Sprintf("%s: not found expected header %s",dest,key))
				return
			}
            if v[0] != value {
				context.JSON(500,fmt.Sprintf("%s: real %s header is %s not equal expect %s",dest,key,v[0],value))
				return
			}
		}
	}

	fmt.Printf("%+v",config)
	context.JSON(200,config)
}