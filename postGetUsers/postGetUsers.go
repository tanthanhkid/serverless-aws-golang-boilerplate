package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"postInsert/models"
	"strconv"
	"time"

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
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId      string        `json:"responseId"`
	ResponseTime    string        `json:"responseTime"`
	ResponseCode    string        `json:"responseCode"`
	ResponseMessage string        `json:"responseMessage"`
	Data            []models.User `json:"data"`
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

	//get user from request
	users, err := getuserdetail(bodyRequest.Data.UserName)

	responseCode := "06"
	if err == nil {
		responseCode = "00"
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseId:      uuid.New().String(),
		ResponseTime:    datetime.Format(time.RFC3339),
		ResponseCode:    responseCode,
		ResponseMessage: "rows: " + strconv.Itoa(len(users)),
		Data:            users,
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

func getuserdetail(username string) ([]models.User, error) {

	//create connection
	db := createConnection()

	//close the connection
	defer db.Close()

	//create user array
	var users []models.User

	//create insert query
	sqlStatement := `SELECT id, username, name, phone FROM users WHERE username = $1`

	rows, err := db.Query(sqlStatement, username)

	if err != nil {
		fmt.Println("system err when exec sql statement:" + err.Error())
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Name, &user.Phone); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	fmt.Println("get users success")

	return users, err

}
