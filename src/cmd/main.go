package main

import (
	"app-gateway/src/proto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	r := gin.Default()

	r.POST("/send", func(c *gin.Context) {
		content := c.PostForm("content")

		conn, err := grpc.Dial(":50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := proto.NewMessageServiceClient(conn)
		res, err := client.SendMessage(c, &proto.SendMessageRequest{Content: content})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": res.Status})
	})

	r.GET("/receive/:id", func(c *gin.Context) {
		id := c.Param("id")

		conn, err := grpc.Dial(":50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := proto.NewMessageServiceClient(conn)
		res, err := client.ReceiveMessage(c, &proto.ReceiveMessageRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"content": res.Content})
	})

	r.Run(":8080")
}
