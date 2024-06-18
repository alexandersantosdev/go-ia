package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Message struct {
	Content string `json:"content"`
}

var IA_KEY string = ""
var PORT string = "8080"
var GIN_MODE string = "debug"

func main() {

	gin.SetMode(gin.DebugMode)

	if GIN_MODE == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	ctx := context.Background()
	gClient, err := genai.NewClient(ctx, option.WithAPIKey(IA_KEY))

	if err != nil {
		panic(err)
	}

	defer gClient.Close()

	model := gClient.GenerativeModel("gemini-1.0-pro")
	cs := model.StartChat()

	send := func(msg string) *genai.GenerateContentResponse {
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/chat", func(c *gin.Context) {
		var msg Message
		c.BindJSON(&msg)
		res := send(msg.Content)
		c.JSON(200, printResponse(res))
	})

	fmt.Println("Listening on port " + PORT)

	if err := router.Run(fmt.Sprintf(":%s", PORT)); err != nil {
		panic(err)
	}
}

func init() {
	godotenv.Load(".env")
	IA_KEY = os.Getenv("IA_KEY")
	PORT = os.Getenv("PORT")
	GIN_MODE = os.Getenv("GIN_MODE")
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var resposta string = ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				resposta = fmt.Sprintf("%s\n%s", resposta, part)
			}
		}
	}
	return resposta
}
