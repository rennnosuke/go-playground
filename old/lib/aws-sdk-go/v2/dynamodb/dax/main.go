package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "test-dax-game-score"

const help = `invalid arguments.
please input command as following:

- get item:  dax-test put 1 1
- put item:  dax-test put 1 1 100
`

func main() {
	args := valueArgs()

	if len(args) <= 1 {
		_, _ = fmt.Fprintf(os.Stderr, help)
		os.Exit(0)
	}

	// dynamodb client
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	switch args[0] {
	case "put":
		if err := putItem(ctx, client, args[1:]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to put item: err=%v", err)
			os.Exit(0)
		}
	case "get":
		item, err := getItem(ctx, client, args[1:])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to get item: err=%v", err)
			os.Exit(0)
		}
		jsonStr, err := asJSONStr(item)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to parse item as json: err=%v", err)
			os.Exit(0)
		}
		fmt.Println(jsonStr)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "invalid arguments.")
		os.Exit(0)
	}
}

func getItem(ctx context.Context, client *dynamodb.Client, args []string) (map[string]types.AttributeValue, error) {
	// validate arguments
	if err := validateArgs(args, 2); err != nil {
		return nil, err
	}

	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberN{Value: args[0]},
			"game_id": &types.AttributeValueMemberN{Value: args[1]},
		},
		TableName: aws.String(tableName),
	}

	res, err := client.GetItem(ctx, &input)
	if err != nil {
		return nil, err
	}

	return res.Item, nil
}

func putItem(ctx context.Context, client *dynamodb.Client, args []string) error {
	// validate arguments
	if err := validateArgs(args, 3); err != nil {
		return err
	}

	// put item
	item := map[string]types.AttributeValue{
		"user_id": &types.AttributeValueMemberN{Value: args[0]},
		"game_id": &types.AttributeValueMemberN{Value: args[1]},
		"score":   &types.AttributeValueMemberN{Value: args[2]},
	}
	input := dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err := client.PutItem(ctx, &input)
	return err
}

func valueArgs() []string {
	if len(os.Args) < 2 {
		return nil
	}
	if os.Args[0] == "go" {
		return os.Args[2:]
	}
	return os.Args[1:]
}

func validateArgs(args []string, argLen int) error {
	if len(args) != argLen {
		return fmt.Errorf("len(args) is not %d, got=%d(%v)", argLen, len(args), args)
	}
	for i, arg := range args {
		if _, err := strconv.Atoi(arg); err != nil {
			return fmt.Errorf("%v: failed to parse argument[%d]: %s", err, i, arg)
		}
	}
	return nil
}

func asJSONStr(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
