package bot

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	TL_API  string = "http://openapi.turingapi.com/openapi/api/v2"
	USER_ID string = "77"
)

type ResText struct {
	resultType string `jsop:"resultType"`
	values     string `jsop:"values.text"`
}

type Intent struct {
	Code int `json:"code"`
}

type Results struct {
	ResultType string            `json:"resultType"`
	Values     map[string]string `json:"values"`
}

type Result struct {
	Status  Intent    `json:"intent"`
	Results []Results `json:"results"`
}

func TlBot(info string, apiKey string) string {
	client := resty.New() // 创建一个restry客户端
	client.SetProxy("http://127.0.0.1:12639")
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
	})
	resp, err := client.R().EnableTrace().SetBody([]byte(
		fmt.Sprintf(`{
			"perception": {
					"inputText": {
							"text": "%v"
					},
			},
			"userInfo": {
					"apiKey": "%v",
					"userId": "%v"
			}
	}`, info, apiKey, USER_ID),
	)).Post(TL_API)
	if err != nil {
		fmt.Printf("请求接口失败：%v\n", err)
		return ""
	}
	var resContent Result
	fmtErr := json.Unmarshal([]byte(resp.String()), &resContent)
	if fmtErr != nil {
		fmt.Printf("解析错误: %v\n", fmtErr)
		return ""
	}
	fmt.Printf("解析结果：%v", resContent)
	if resContent.Status.Code != 0 {
		fmt.Printf("请求接口失败：%v\n", resContent.Status.Code)
		return ""
	}
	return formatValues(&resContent)
}

func formatValues(res *Result) string {
	result := ""
	for _, value := range res.Results {
		if value.ResultType == "text" {
			for _, v := range value.Values {
				result += v
			}
		}
	}
	return result
}
