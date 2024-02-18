// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

/*
*

  - 使用AK&SK初始化账号Client

  - @param accessKeyId

  - @return Client

  - @throws Exception
*/
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *ecs20140526.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Ecs
	config.Endpoint = tea.String("ecs.cn-shanghai.aliyuncs.com")
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	if len(args) == 0 {
		fmt.Println("Usage: update_aliyun_firewall_for_frpc <location tag>")
		os.Exit(1)
	}
	locationTag := *args[0]
	client, _err := CreateClient(tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")), tea.String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")))
	if _err != nil {
		return _err
	}

	// rules := list.New()
	var rules []string
	var publicIpAddr string

	// get public IP address of my compute
	requestURL := "https://www.ipplus360.com/getIP"
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	var data interface{}
	d := json.NewDecoder(strings.NewReader(string(resBody)))
	d.Decode(&data)
	if m, ok := data.(map[string]interface{}); ok {
		if ipAddr, found := m["data"]; found {
			publicIpAddr = ipAddr.(string)
			fmt.Printf("client: got public ip address: %s\n", publicIpAddr)
		} else {
			fmt.Printf("client: can't find data field in response json payload.\n")
			os.Exit(1)
		}
	} else {
		fmt.Printf("client: failed to convert body string to map[string]interface{}. Err: %t\n", ok)
		os.Exit(1)
	}

	// get rules with given location in Description field
	describeSecurityGroupAttributeRequest := &ecs20140526.DescribeSecurityGroupAttributeRequest{
		RegionId:        tea.String("cn-shanghai"),
		SecurityGroupId: tea.String("sg-uf69xcfy8qawo9wj6vot"),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, _err := client.DescribeSecurityGroupAttributeWithOptions(describeSecurityGroupAttributeRequest, runtime)
		if _err != nil {
			return _err
		}

		// fmt.Printf("%v", resp)

		for _, item := range resp.Body.Permissions.Permission {
			// fmt.Printf("%v", item)
			ruleId := item.SecurityGroupRuleId
			descrption := item.Description
			// if strings.LastIndex(*description, )
			if strings.LastIndex(*descrption, locationTag) >= 0 {
				// fmt.Printf("ruleId: %s\n", *ruleId)
				rules = append(rules, *ruleId)
			}
		}

		// modify the rules to change the source ip
		for _, ruleId := range rules {

			fmt.Printf("aliyun client: changing SourceCidrIp to %v for rule %s (desc includes '%s').\n", publicIpAddr, ruleId, locationTag)
			modifySecurityGroupRuleRequest := &ecs20140526.ModifySecurityGroupRuleRequest{
				RegionId:            tea.String("cn-shanghai"),
				SecurityGroupId:     tea.String("sg-uf69xcfy8qawo9wj6vot"),
				SourceCidrIp:        tea.String(publicIpAddr),
				SecurityGroupRuleId: tea.String(ruleId),
			}

			_, _err := client.ModifySecurityGroupRuleWithOptions(modifySecurityGroupRuleRequest, runtime)
			// fmt.Printf("resp: %v\n", resp)
			if _err != nil {
				return _err
			}

			fmt.Println(" => OK")
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(error.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
