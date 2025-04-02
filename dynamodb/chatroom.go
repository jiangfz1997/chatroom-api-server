package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

var ChatroomTableName = "chatrooms" // 可换成环境变量或配置文件读取

type Chatroom struct {
	RoomID    string   `json:"room_id" dynamodbav:"room_id"`
	Name      string   `json:"name" dynamodbav:"name"`
	IsPrivate bool     `json:"is_private" dynamodbav:"is_private"`
	CreatedBy string   `json:"created_by" dynamodbav:"created_by"`
	CreatedAt string   `json:"created_at" dynamodbav:"created_at"`
	Users     []string `json:"users" dynamodbav:"users"`
}

func CreateChatroomTable() {
	_, err := DB.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(ChatroomTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("room_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("room_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("❌ 创建 chatrooms 表失败: %v", err)
	}
	log.Println("✅ 聊天室表创建成功")
}
func CreateChatroom(chatroom Chatroom) error {
	// 时间格式化（标准 ISO 时间）
	if chatroom.CreatedAt == "" {
		chatroom.CreatedAt = time.Now().Format(time.RFC3339)
	}

	item, err := attributevalue.MarshalMap(chatroom)
	if err != nil {
		return err
	}

	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &ChatroomTableName, // 推荐把表名用变量保存
		Item:      item,
	})
	return err
}

func GetChatroom(chatroomId string) (Chatroom, error) {
	var chatroom Chatroom

	// 查询条件
	input := &dynamodb.GetItemInput{
		TableName: aws.String(ChatroomTableName),
		Key: map[string]types.AttributeValue{
			"room_id": &types.AttributeValueMemberS{Value: chatroomId},
		},
	}

	result, err := DB.GetItem(context.TODO(), input)
	if err != nil {
		return chatroom, err
	}

	if result.Item == nil {
		return chatroom, fmt.Errorf("聊天室不存在")
	}

	err = attributevalue.UnmarshalMap(result.Item, &chatroom)
	if err != nil {
		return chatroom, err
	}

	return chatroom, nil
}

func AddUserToChatroom(username, roomID string) error {
	// Step 1: 先获取聊天室对象
	chatroom, err := GetChatroom(roomID)
	if err != nil {
		return fmt.Errorf("聊天室不存在: %w", err)
	}

	// Step 2: 检查用户是否已经存在
	for _, u := range chatroom.Users {
		if u == username {
			return nil // 已存在就跳过
		}
	}

	// Step 3: 添加新用户
	chatroom.Users = append(chatroom.Users, username)

	// Step 4: 写回数据库
	item, err := attributevalue.MarshalMap(chatroom)
	if err != nil {
		return err
	}
	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(ChatroomTableName),
		Item:      item,
	})
	return err
}

func RemoveUserFromChatroom(username, roomID string) error {
	room, err := GetChatroom(roomID)
	if err != nil {
		return fmt.Errorf("聊天室不存在: %w", err)
	}

	// 过滤掉这个用户
	var newUsers []string
	for _, u := range room.Users {
		if u != username {
			newUsers = append(newUsers, u)
		}
	}
	room.Users = newUsers

	// 写回数据库
	item, err := attributevalue.MarshalMap(room)
	if err != nil {
		return err
	}
	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(ChatroomTableName),
		Item:      item,
	})
	return err
}

func GetChatroomsByUsername(username string) ([]Chatroom, error) {
	var results []Chatroom

	// 全表扫描
	output, err := DB.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(ChatroomTableName),
	})
	if err != nil {
		return nil, err
	}

	for _, item := range output.Items {
		var room Chatroom
		if err := attributevalue.UnmarshalMap(item, &room); err != nil {
			continue
		}

		// 判断 users 中是否包含该用户
		for _, u := range room.Users {
			if u == username {
				results = append(results, room)
				break
			}
		}
	}

	return results, nil
}
