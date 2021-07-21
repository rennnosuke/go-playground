package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "test-dax-game-score"

func main() {
	if len(os.Args) <= 1 {
		_, _ = fmt.Fprintf(os.Stderr, "invalid arguments.\n")
		_, _ = fmt.Fprintf(os.Stderr, "please input command as following:\n")
		_, _ = fmt.Fprintf(os.Stderr, "- get item:  dax-test put 1 1.\n")
		_, _ = fmt.Fprintf(os.Stderr, "- put item:  dax-test put 1 1 100.\n")
		os.Exit(0)
	}

	// dynamodb client
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	client := dynamodb.NewFromConfig(cfg)

	switch os.Args[1] {
	case "put":
		if err := putItem(ctx, client, valueArgs()[1:]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to put item: err=%v", err)
			os.Exit(0)
		}
	case "get":
	default:
		_, _ = fmt.Fprintf(os.Stderr, "invalid arguments.")
		os.Exit(0)
	}
}

func putItem(ctx context.Context, client *dynamodb.Client, args []string) error {
	// validate arguments
	const argLen = 3
	if len(args) != argLen {
		return fmt.Errorf("len(args) is not %d, got=%d(%v)", argLen, len(args), args)
	}
	for i, arg := range args {
		if _, err := strconv.Atoi(arg); err != nil {
			return fmt.Errorf("%v: failed to parse argument[%d]: %s", err, i, arg)
		}
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
