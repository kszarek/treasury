package ssm

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	Key1 = "test/webapp/cocpit_api_pass"
	Key2 = "test/webapp/user_api_pass"
	Key3 = "test/cockpit/user_api_pass"
)

var KeyValueMap = map[string]string{
	Key1: "as9@#$%^&*(/2hdiwnf",
	Key2: "as9@#$&*(/2saddsahdiwnf",
	Key3: "#$&*(/2saddsah&as",
}

const (
	SSMKey1 = "/" + Key1
	SSMKey2 = "/" + Key2
	SSMKey3 = "/" + Key3
)

var SSMKeyValueMap = map[string]string{
	SSMKey1: KeyValueMap[Key1],
	SSMKey2: KeyValueMap[Key2],
	SSMKey3: KeyValueMap[Key3],
}

// MockSSMClient fake SSMAPI
type MockSSMClient struct {
	ssmiface.SSMAPI
}

func (m *MockSSMClient) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("PutParameterInput is empty")
	}
	if *input.Name == "" {
		return nil, fmt.Errorf("Name in PutParameterInput is not set")
	}
	// SSM path based key needs to start from "/"
	name := *input.Name
	if string(name[0]) != "/" {
		return nil, fmt.Errorf("SSM Name needs to start from /")
	}
	if *input.Value == "" {
		return nil, fmt.Errorf("Value in PutParameterInput is not set")
	}
	var version int64 = 1
	return &ssm.PutParameterOutput{Version: &version}, nil
}

func (m *MockSSMClient) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	log.Println("input.Name:", input.Name)
	if _, ok := SSMKeyValueMap[*input.Name]; !ok {
		return nil, fmt.Errorf("Missing key:%s in KeyValue map", *input.Name)
	}
	if !*input.WithDecryption {
		return nil, fmt.Errorf("Missing decryption field")
	}
	value := SSMKeyValueMap[*input.Name]
	return &ssm.GetParameterOutput{
		// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#Parameter
		Parameter: &ssm.Parameter{
			Name:  input.Name,
			Value: &value,
		},
	}, nil
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#SSM.GetParametersByPath
// https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/#GetParametersByPathInput
func (m *MockSSMClient) GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	if !*input.WithDecryption {
		return nil, fmt.Errorf("Missing decryption field")
	}
	var contents []*ssm.Parameter
	for key, value := range SSMKeyValueMap {
		key := key
		value := value
		if strings.Contains(key, *input.Path) {
			contents = append(contents, &ssm.Parameter{
				Name:  &key,
				Value: &value,
			})
		}
	}
	return &ssm.GetParametersByPathOutput{
		Parameters: contents,
	}, nil
}
