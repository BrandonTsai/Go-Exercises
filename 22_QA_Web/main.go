package main

import (
	"gitbub.com/brandontsai/qa/api"
	"gitbub.com/brandontsai/qa/common"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitDB()

	r := gin.Default()
	v1 := r.Group("/api")
	api.LabelsRegister(v1)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
