package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const labelURL = "/labels"

func LabelsRegister(r *gin.RouterGroup) {
	r.GET(labelURL, GetLabels)
	r.POST(labelURL, CreateLabels)
	r.DELETE(labelURL+"/:id", DeleteLabels)
	r.PUT(labelURL+"/:id", UpdateLabels)
}

func GetLabels(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateLabels(c *gin.Context) {
	fmt.Println("Create labels")
}

func DeleteLabels(c *gin.Context) {
	fmt.Println("Delete labels")
}

func UpdateLabels(c *gin.Context) {
	fmt.Println("Update labels")
}
