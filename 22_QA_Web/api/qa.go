package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const qaURL = "/qa"

func QARegister(r *gin.RouterGroup) {
	r.GET(qaURL, GetQAs)
	r.POST(qaURL, CreateQAs)
	r.DELETE(qaURL+"/:id", DeleteQAs)
	r.PUT(qaURL+"/:id", UpdateQAs)
}

func GetQAs(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get  questions",
	})
}

func CreateQAs(c *gin.Context) {
	fmt.Println("Create QAs")
}

func DeleteQAs(c *gin.Context) {
	fmt.Println("Delete QAs")
}

func UpdateQAs(c *gin.Context) {
	fmt.Println("Update QAs")
}
