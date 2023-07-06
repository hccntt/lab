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
	RequestId   string  `json:"requestId"`
	RequestTime string  `json:"requestTime"`
	Data        DataReq `json:"data"`
}

type DataReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
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

type ResponseDataApi2 struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type ResponseApi2 struct {
	ResponseId   string           `json:"responseId"`
	ResponseTime string           `json:"responseTime"`
	Data         ResponseDataApi2 `json:"data"`
}

// Handler function Using AWS Lambda Proxy Request
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	datetime := time.Now().UTC()
	// BodyRequest will be used to take the json response from client and build it
	bodyRequest := BodyRequest{}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		errrr := fmt.Sprintf("error json bodyRequest: %s", err.Error())
		return events.APIGatewayProxyResponse{Body: errrr, StatusCode: 400}, nil
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
	// if bodyRequest.Data.Value == nil {
	// 	return events.APIGatewayProxyResponse{Body: "Value1, Value2 can not be null", StatusCode: 400}, nil
	// }

	strData := fmt.Sprintf(`{
		"requestId": "%s",		
		"data": {
			"value": %s
		}
	}`, bodyRequest.RequestId, bodyRequest.Data.Phone) // "requestTime": "%s",   bodyRequest.RequestTime,

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
		errrr := fmt.Sprintf("error json dataApi: %s", err3.Error())
		return events.APIGatewayProxyResponse{Body: errrr, StatusCode: 400}, nil
	}

	if dataApi.ResponseCode != "00" {
		return events.APIGatewayProxyResponse{Body: "Do not allow insert data", StatusCode: 400}, nil
	}

	// call api insert data

	strData2 := fmt.Sprintf(`{
		"requestId": "%s",
		"requestTime": "%s",
		"data": {
			"username": "%s",
			"name": "%s",
			"phone": "%s"
		}
	}`, bodyRequest.RequestId, bodyRequest.RequestTime, bodyRequest.Data.Username, bodyRequest.Data.Name, bodyRequest.Data.Phone)
	//var jsonData = []byte(strData2)

	url2 := "https://q7nyrvbdjb.execute-api.us-east-1.amazonaws.com/dev/postapi3"
	//contentType2 := "text/plain"
	data2 := []byte(strData2)

	client2 := &http.Client{}
	req2, err2 := http.NewRequest("POST", url2, bytes.NewBuffer(data2))
	if err2 != nil {
		fmt.Println(err2)
		//return
	}
	//req.Header.Add("Content-Type", contentType2)
	req.Header.Add("x-api-key", "ezBKpxf3w24tY5dVSBaap6O42ZIHAQlW3IhK5ZwF")

	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		fmt.Println(err2)
		//return
	}
	defer resp2.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println(err)
		//return
	}

	dataApi2 := ResponseDataApi2{}
	err4 := json.Unmarshal(body2, &dataApi2)

	if err4 != nil {
		errrr := fmt.Sprintf("error json dataApi2: %s - data: %s", err4.Error(), strData2)
		return events.APIGatewayProxyResponse{Body: errrr, StatusCode: 400}, nil
	}

	bodyResponse := ResponseApi2{
		ResponseId:   uuid.New().String(),
		ResponseTime: datetime.Format(time.RFC3339),
		Data:         ResponseDataApi2(dataApi2),
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
