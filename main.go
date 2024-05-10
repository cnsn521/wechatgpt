package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 请求数据的结构体
type RequestData struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

// 响应数据的结构体
type ResponseData struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
	Usage  struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

// 定义 JSON 数据的结构体
type Payload struct {
	AdminIdList []string `json:"adminIdList"`
	Avatar      string   `json:"avatar"`
	ID          string   `json:"id"`
	Topic       string   `json:"topic"`
	MemberList  []struct {
		Avatar string `json:"avatar"`
		ID     string `json:"id"`
		Name   string `json:"name"`
		Alias  string `json:"alias"`
	} `json:"memberList"`
}

type RoomData struct {
	Room struct {
		Payload Payload `json:"payload"`
	} `json:"room"`
	From struct {
		Payload struct {
			Name string `json:"name"`
			Type int    `json:"type"`
		} `json:"payload"`
	} `json:"from"`
}

// 接收的消息
type Message struct {
	Type          string `json:"type"`
	Content       string `json:"content"`
	Source        string `json:"source"`
	IsMentioned   string `json:"isMentioned"`
	IsMsgFromSelf string `json:"isMsgFromSelf"`
	IsSystemEvent string `json:"isSystemEvent"`
}

// 发送的消息内容
type Data struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// 发送的消息
type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

func main() {
	http.HandleFunc("/msg", handleRequest)
	fmt.Println("Starting server on port 3002")
	http.ListenAndServe(":3002", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 解析多段表单数据
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	var message Message
	// 迭代处理每个表单部分
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, "Error parsing form part", http.StatusBadRequest)
			return
		}

		// 处理表单部分
		switch part.FormName() {
		case "type":
			// 处理 "type" 部分的内容
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.Type = string(data)
			fmt.Println("Type:", string(data))
		case "content":
			// 处理 "content" 部分的内容
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.Content = string(data)
			fmt.Println("Content:", string(data))
		case "source":
			// 处理 "source" 部分的内容
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.Source = string(data)
			fmt.Println("Source:", string(data))
		case "isMentioned":
			// 处理 "source" 部分的内容,该消息是@我的消息
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.IsMentioned = string(data)
			fmt.Println("isMentioned:", string(data))
		case "isMsgFromSelf":
			// 处理 "source" 部分的内容,是否是来自自己的消息
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.IsMsgFromSelf = string(data)
			fmt.Println("isMsgFromSelf:", string(data))
		case "isSystemEvent":
			// 处理 "source" 部分的内容,是否是来自自己的消息
			data, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "Error reading form part data", http.StatusBadRequest)
				return
			}
			message.IsSystemEvent = string(data)
			// fmt.Println("isSystemEvent:", string(data))
			if message.IsSystemEvent == "1" {
				return
			}
		default:
			// 处理其他表单部分
			fmt.Println("Unknown form part:", part.FormName())
		}
	}

	var data RoomData
	if err := json.Unmarshal([]byte(message.Source), &data); err != nil {
		panic(err)
	}
	var response Response
	if message.Type != "text" {
		response.Success = true
		response.Data.Type = "text"
		response.Data.Content = "只支持文本，请重新输入文本!"
	} else {
		// 判断 room 是否为空
		if data.Room.Payload.Topic == "" {
			// 如果为空，则输出 from 中的 name 信息
			fmt.Println("Room is empty.")
			fmt.Println("Name:", data.From.Payload.Name)
			if data.From.Payload.Type == 1 {
				response.Success = true
				response.Data.Type = "text"
				response.Data.Content = aurora(message.Content)
			} else {
				return
			}
		} else {
			// 如果不为空，则输出 topic 的名称和 from 中的 name 信息
			fmt.Println("Topic:", data.Room.Payload.Topic)
			fmt.Println("Name:", data.From.Payload.Name)
			response.Success = true
			if message.IsMentioned == "1" {
				response.Data.Type = "text"
				response.Data.Content = "@" + data.From.Payload.Name + " " + aurora(message.Content)
			} else {

				return
			}
		}
	}
	// 编码响应为 JSON 格式
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(jsonData))
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应
	w.Write(jsonData)

}

func aurora(str string) string {
	var msg_content string

	// 构造请求数据
	requestData := RequestData{
		Model: "gpt-3.5-turbo",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: str,
			},
		},
		Stream: false,
	}

	// 将请求数据转换为 JSON 格式
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		panic(err)
	}

	// 发送 POST 请求
	resp, err := http.Post("https://aurora.xncen.top/v1/chat/completions", "application/json", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	responseDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析响应数据
	var responseData ResponseData
	if err := json.Unmarshal(responseDataBytes, &responseData); err != nil {
		panic(err)
	}
	msg_content = responseData.Choices[0].Message.Content
	// 输出 content 内容
	fmt.Println("Content:", msg_content)

	return msg_content
}
