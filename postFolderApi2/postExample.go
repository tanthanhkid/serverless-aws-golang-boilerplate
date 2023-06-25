package main

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	"encoding/hex"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestId   string      `json:"requestId"`
	RequestTime string      `json:"requestTime"`
	Data        DataRequest `json:"data"`
}

type DataRequest struct {
	PlainText string `json:"plainText"`
	SecretKey string `json:"secretKey"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId   string       `json:"responseId"`
	ResponseTime string       `json:"responseTime"`
	Data         DataResponse `json:"data"`
}

type DataResponse struct {
	Signature string `json:"signature"`
}

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	datetime := time.Now().UTC()
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{
		RequestId: "",
	}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 401}, nil
	}

	//verify uuid not null
	if bodyRequest.RequestId == "" {
		return events.APIGatewayProxyResponse{Body: "requestId can not be null", StatusCode: 401}, nil

	}

	//verify datetime format RFC3339
	parsedTime, err := time.Parse(time.RFC3339, bodyRequest.RequestTime)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error() + "parsedTime: " + parsedTime.GoString(), StatusCode: 401}, nil
	}

	//verify hash materials
	if bodyRequest.Data.PlainText == "" || bodyRequest.Data.SecretKey == "" {
		return events.APIGatewayProxyResponse{Body: "PlaintText, SecretKey can not be null", StatusCode: 401}, nil

	}

	//do the hashing works
	h := sha256.New()
	h.Write([]byte(bodyRequest.Data.PlainText + bodyRequest.Data.SecretKey))
	sha256_hash := hex.EncodeToString(h.Sum(nil))

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseId:   uuid.New().String(),
		ResponseTime: datetime.Format(time.RFC3339),
		Data:         DataResponse{Signature: sha256_hash},
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
