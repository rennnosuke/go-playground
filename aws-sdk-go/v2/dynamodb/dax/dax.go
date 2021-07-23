package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Command string

const GET Command = "get"
const PUT Command = "put"

type HandlerParameters struct {
	Endpoint  string
	Region    string
	TableName string
	Cmd       Command
	Args      map[string]interface{}
}

func handleDax(input HandlerParameters) (map[string]interface{}, error) {
	client, err := connect(input.Endpoint, input.Region)
	if err != nil {
		return nil, err
	}

	switch input.Cmd {
	case "get":
		return getDaxItem(client, input.TableName, input.Args)
	case "put":
	default:
	}

	return nil, nil
}

func getDaxItem(client *dax.Dax, tableName string, keys map[string]interface{}) (map[string]interface{}, error) {
	kv := make(map[string]*dynamodb.AttributeValue, len(keys))
	for k, v := range keys {
		kv[k] = attributeValue(v)
	}

	input := dynamodb.GetItemInput{
		Key:       kv,
		TableName: aws.String(tableName),
	}
	out, err := client.GetItem(&input)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{}, len(out.Item))
	for k, v := range out.Item {
		switch {
		case v.BOOL != nil:
			res[k] = v.BOOL
		case v.S != nil:
			res[k] = v.S
		case v.N != nil:
			res[k] = v.N
		case v.M != nil:
			res[k] = v.M
		case v.B != nil:
			res[k] = v.B
		}
	}

	return res, nil
}

func attributeValue(x interface{}) *dynamodb.AttributeValue {
	switch v := x.(type) {
	case int, int32, int64:
		return &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%d", v)),
		}
	case string:
		return &dynamodb.AttributeValue{
			S: aws.String(v),
		}
	case bool:
		return &dynamodb.AttributeValue{
			BOOL: aws.Bool(v),
		}
	}
	return nil
}

func connect(endpoint string, region string) (*dax.Dax, error) {
	if strings.HasPrefix(endpoint, "dax://") {
		return connectCluster(endpoint, region)
	}
	if strings.HasPrefix(endpoint, "daxs://") {
		return connectSecureCluster(endpoint, region)
	}
	return nil, fmt.Errorf("invalid endpoint: %s", endpoint)
}

func connectCluster(endpoint string, region string) (*dax.Dax, error) {
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{endpoint}
	cfg.Region = region
	return dax.New(cfg)
}

func connectSecureCluster(endpoint string, region string) (*dax.Dax, error) {
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{endpoint}
	cfg.Region = region
	cfg.SkipHostnameVerification = false
	cfg.DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
		dialCon, err := dax.SecureDialContext(endpoint, cfg.SkipHostnameVerification)
		if err != nil {
			return nil, err
		}
		return dialCon(ctx, network, address)
	}
	return dax.New(cfg)
}
