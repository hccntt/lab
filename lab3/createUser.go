package main

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BodyRequest is our self-made struct to process JSON request from Client
type BodyRequest struct {
	RequestId   string      `json:"requestId"`
	RequestTime string      `json:"requestTime"`
	Data        DataRequest `json:"data"`
}

type DataRequest struct {
	Value *int `json:"value"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseId   string       `json:"responseId"`
	ResponseTime string       `json:"responseTime"`
	Data         DataResponse `json:"data"`
}

type DataResponse struct {
	Response string `json:"response"`
}

type ResponseDataApi struct {
	ResponseId      string `json:"responseId"`
	ResponseTime    string `json:"responseTime"`
	ResponseMessage string `json:"responseMessage"`
	ResponseCode    string `json:"responseCode"`
}

type ResponseApi struct {
	ResponseId   string          `json:"responseId"`
	ResponseTime string          `json:"responseTime"`
	Data         ResponseDataApi `json:"data"`
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
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	//verify uuid not null
	if bodyRequest.RequestId == "" {
		return events.APIGatewayProxyResponse{Body: "requestId can not be null", StatusCode: 400}, nil
	}

	//verify datetime format RFC3339
	parsedTime, err := time.Parse(time.RFC3339, bodyRequest.RequestTime)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error() + "parsedTime: " + parsedTime.GoString(), StatusCode: 400}, nil
	}

	//verify sum materials
	if bodyRequest.Data.Value == nil {
		return events.APIGatewayProxyResponse{Body: "Value1, Value2 can not be null", StatusCode: 400}, nil
	}

	strData := fmt.Sprintf(`{
		"requestId":"requestId",
		"data": {
			"value": %d
		}
	}`, bodyRequest.Data.Value)

	url := "https://1g1zcrwqhj.execute-api.ap-southeast-1.amazonaws.com/dev/testapi"
	contentType := "text/plain"
	data := []byte(strData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		//return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-api-key", "B5d4JtTU8u1ggV8gp7OF88gcCGxZls6T3f5PYZSa")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		msg := fmt.Sprintf("Non-OK HTTP status: %d", resp.StatusCode)
		// You may read / inspect response body
		return events.APIGatewayProxyResponse{Body: msg, StatusCode: 400}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		//return
	}

	//fmt.Println(string(body))

	// We will build the BodyResponse and send it back in json form
	// bodyResponse := BodyResponse{
	// 	ResponseId:   uuid.New().String(),
	// 	ResponseTime: datetime.Format(time.RFC3339),
	// 	Data:         DataResponse{Response: string(body)},
	// }
	dataApi := ResponseDataApi{}
	err3 := json.Unmarshal(body, &dataApi)

	if err3 != nil {
		return events.APIGatewayProxyResponse{Body: err3.Error(), StatusCode: 400}, nil
	}

	bodyResponse := ResponseApi{
		ResponseId:   uuid.New().String(),
		ResponseTime: datetime.Format(time.RFC3339),
		Data:         dataApi,
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
