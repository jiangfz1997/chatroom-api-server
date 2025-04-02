package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

import "time"

var MessageTableName = "messages"

type Message struct {
	RoomID    string `json:"room_id" dynamodbav:"room_id"`     // 分区键
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"` // 排序键
	Sender    string `json:"sender" dynamodbav:"sender"`
	Text      string `json:"text" dynamodbav:"text"`
}

func NewMessage(roomID, sender, text string) Message {
	return Message{
		RoomID:    roomID,
		Sender:    sender,
		Text:      text,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func GetMessagesBefore(roomID, before string, limit int) ([]Message, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(MessageTableName),
		KeyConditionExpression: aws.String("room_id = :rid AND #ts < :before"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp", // DynamoDB 的关键字需要映射
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":rid":    &types.AttributeValueMemberS{Value: roomID},
			":before": &types.AttributeValueMemberS{Value: before},
		},
		Limit:            aws.Int32(int32(limit)),
		ScanIndexForward: aws.Bool(false), // 👈 倒序
	}

	resp, err := DB.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var msgs []Message
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &msgs)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func CreateMessageTable() {
	_, err := DB.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(MessageTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("room_id"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("timestamp"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("room_id"), KeyType: types.KeyTypeHash},    // 分区键
			{AttributeName: aws.String("timestamp"), KeyType: types.KeyTypeRange}, // 排序键
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("❌ 创建 messages 表失败: %v", err)
	}
	log.Println("✅ 消息表创建成功（主键为 room_id + timestamp）")
}
