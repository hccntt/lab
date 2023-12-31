package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"database/sql"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestId   string      `json:"requestId"`
	RequestTime string      `json:"requestTime"`
	Data        DataRequest `json:"data"`
}

type DataRequest struct {
	Value1 *int `json:"value1"`
	Value2 *int `json:"value2"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId   string       `json:"responseId"`
	ResponseTime string       `json:"responseTime"`
	Data         DataResponse `json:"data"`
}

type DataResponse struct {
	Sum int `json:"sum"`
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

	db, errDb := sql.Open("mysql", "hccntt:hccntt123456@tcp(85.10.205.173:3306)/mysqlfree?charset=utf8mb4&parseTime=True&loc=Local") // user:password@tcp(db-hostname:3306)/mydb -- hccntt:hccntt123456@tcp(85.10.205.173:3306)/mysqlfree?charset=utf8mb4&parseTime=True&loc=Local
	if errDb != nil {
		//panic(err.Error())
		return events.APIGatewayProxyResponse{Body: errDb.Error(), StatusCode: 401}, nil
	}
	defer db.Close()

	query := "INSERT INTO `users` (`username`, `name`, `phone`) VALUES (?, ?, ?)"
	insertResult, errI := db.ExecContext(context.Background(), query, "John", "Doe", "33336879879")

	if errI != nil {
		log.Fatalf("impossible insert users: %s", errI)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("impossible insert: %s", errI), StatusCode: 401}, nil
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("impossible insert: %s", err), StatusCode: 401}, nil
	}
	log.Printf("inserted id: %d", id)

	//verify datetime format RFC3339
	parsedTime, err := time.Parse(time.RFC3339, bodyRequest.RequestTime)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error() + "parsedTime: " + parsedTime.GoString(), StatusCode: 401}, nil
	}

	//verify sum materials
	if bodyRequest.Data.Value1 == nil || bodyRequest.Data.Value2 == nil {
		return events.APIGatewayProxyResponse{Body: "Value1, Value2 can not be null", StatusCode: 401}, nil
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseId:   uuid.New().String(),
		ResponseTime: datetime.Format(time.RFC3339),
		Data:         DataResponse{Sum: *bodyRequest.Data.Value1 + *bodyRequest.Data.Value2},
		//Data: DataResponse{Sum: id},
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
