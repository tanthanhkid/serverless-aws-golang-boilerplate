package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	models "postInsert/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestId   string      `json:"requestId"`
	RequestTime string      `json:"requestTime"`
	Data        DataRequest `json:"data"`
}

type DataRequest struct {
	UserName string `json:"userName"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId      string `json:"responseId"`
	ResponseTime    string `json:"responseTime"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
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

	if bodyRequest.Data.Name == "" || bodyRequest.Data.UserName == "" || bodyRequest.Data.Phone == "" {
		return events.APIGatewayProxyResponse{Body: "User object can not be null", StatusCode: 401}, nil
	}

	// insert to database
	var user models.User
	user.Name = bodyRequest.Data.Name
	user.UserName = bodyRequest.Data.UserName
	user.Phone = bodyRequest.Data.Phone

	id := insertUser(user)

	responseCode := "06"
	if id > 0 {
		responseCode = "00"
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseId:      uuid.New().String(),
		ResponseTime:    datetime.Format(time.RFC3339),
		ResponseCode:    responseCode,
		ResponseMessage: "new id is: " + strconv.FormatInt(id, 10),
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
	// Handler(events.APIGatewayProxyRequest{})
	// fmt.Println("ok")
}

func createConnection() *sql.DB {
	connStr := os.Getenv("RDS_CONN_STRING")

	fmt.Println("connection string: " + connStr)

	// Open the connection
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Open connection - err: " + err.Error())
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		fmt.Println("Ping database connection - err: " + err.Error())
		panic(err)
	}

	fmt.Println("connected to database!")

	return db
}

func insertUser(user models.User) int64 {

	//create connection
	db := createConnection()

	//close the connection
	defer db.Close()

	//create insert query
	sqlStatement := `INSERT INTO users (name, username, phone) VALUES ($1, $2, $3) RETURNING id`

	// the inserted id will store in this id
	var id int64

	err := db.QueryRow(sqlStatement, user.Name, user.UserName, user.Phone).Scan(&id)

	if err != nil {
		fmt.Println("system err:" + err.Error())
	}
	if id > 0 {
		fmt.Println("insert user success!")
	}

	return id
}
