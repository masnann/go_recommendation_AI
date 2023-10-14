package main

import (
	_ "encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	e := echo.New()
	e.POST("/recommend-laptop", recommendLaptop)
	e.Logger.Fatal(e.Start(":8080"))
}

type LaptopRecommendationRequest struct {
	Budget     int    `json:"budget"`
	Purpose    string `json:"purpose"`
	Brand      string `json:"brand"`
	RAM        string `json:"ram"`
	CPU        string `json:"cpu"`
	ScreenSize string `json:"screen_size"`
}

type LaptopRecommendationResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func recommendLaptop(c echo.Context) error {
	req := new(LaptopRecommendationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Status Bad Request",
		})
	}

	chatMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf("Recommend a laptop for %s use with a budget of %d Rupiah, %s RAM, %s CPU, %s screen size and %s Brand", req.Purpose, req.Budget, req.RAM, req.CPU, req.ScreenSize, req.Brand),
	}

	apiKey := os.Getenv("APIKEY")

	client := openai.NewClient(apiKey)
	chatReq := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{chatMessage},
	}

	resp, err := client.CreateChatCompletion(c.Request().Context(), chatReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Status Internal Server Error",
		})
	}

	recommendation := resp.Choices[0].Message.Content
	result := LaptopRecommendationResponse{
		Status: "success",
		Data:   recommendation,
	}
	return c.JSON(http.StatusOK, result)
}
