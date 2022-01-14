package main

import "github.com/gin-gonic/gin"

func main(){
	server := gin.Default()
	_ = server.Run(":8090")
}

