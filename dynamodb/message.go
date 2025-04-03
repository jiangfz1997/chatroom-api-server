package dynamodb

import (
	log "chatroom-api/logger"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	log.Log.Infof("查询历史消息: room=%s, before=%s, limit=%d", roomID, before, limit)
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
		ScanIndexForward: aws.Bool(false), // 倒序
	}

	resp, err := DB.Query(context.TODO(), input)
	if err != nil {
		log.Log.Errorf("查询消息失败: %v", err)
		return nil, err
	}

	var msgs []Message
	err = attributevalue.UnmarshalListOfMaps(resp.Items, &msgs)
	if err != nil {
		log.Log.Errorf("消息反序列化失败: %v", err)
		return nil, err
	}

	log.Log.Infof("成功查询到 %d 条消息", len(msgs))
	return msgs, nil
}

func CreateMessageTable() {
	log.Log.Info("开始创建 messages 表")
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
		log.Log.Fatalf("创建 messages 表失败: %v", err)
	}

	log.Log.Info("messages 表创建成功（主键为 room_id + timestamp）")
}
