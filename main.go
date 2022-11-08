package main

import (
	"os"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

    g.GET("/ping",handlePing)

    g.GET("/dial",handleCus)

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

func handleCus(context *gin.Context) {
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

	dispatch := CusDispatch{Config: &config, Context: context}

	dispatch.Dispatch()
}

func handlePing(context *gin.Context) {
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

	dispatch := PingDispatch{Config: &config, Context: context}

	dispatch.Dispatch()
}