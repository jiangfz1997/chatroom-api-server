package dynamodb

import (
	log "chatroom-api/logger"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	//"log"
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
	log.Log.Info("准备创建 chatrooms 表")
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
		log.Log.Fatalf("创建 chatrooms 表失败: %v", err)
	}
	log.Log.Info("chatrooms 表创建成功")
}
func CreateChatroom(chatroom Chatroom) error {
	// 时间格式化（标准 ISO 时间）
	if chatroom.CreatedAt == "" {
		chatroom.CreatedAt = time.Now().Format(time.RFC3339)
		log.Log.Debugf("聊天室时间未指定，已自动填充为: %s", chatroom.CreatedAt)
	}
	log.Log.Infof("准备创建聊天室: room_id=%s, name=%s, created_by=%s", chatroom.RoomID, chatroom.Name, chatroom.CreatedBy)
	item, err := attributevalue.MarshalMap(chatroom)
	if err != nil {
		log.Log.Errorf("聊天室数据序列化失败: %v", err)
		return err
	}

	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &ChatroomTableName, // 推荐把表名用变量保存
		Item:      item,
	})
	if err != nil {
		log.Log.Errorf("写入聊天室失败: %v", err)
	} else {
		log.Log.Infof("聊天室创建成功: room_id=%s", chatroom.RoomID)
	}
	return err
}

func GetChatroom(chatroomId string) (Chatroom, error) {
	var chatroom Chatroom
	log.Log.Infof("尝试获取聊天室: room_id=%s", chatroomId)
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
		log.Log.Warnf("未找到聊天室: room_id=%s", chatroomId)
		return chatroom, fmt.Errorf("聊天室不存在")
	}

	err = attributevalue.UnmarshalMap(result.Item, &chatroom)
	if err != nil {
		log.Log.Errorf("聊天室反序列化失败: %v", err)
		return chatroom, err
	}

	log.Log.Infof("成功获取聊天室: room_id=%s", chatroomId)
	return chatroom, nil
}

func AddUserToChatroom(username, roomID string) error {
	log.Log.Infof("尝试将用户加入聊天室: user=%s, room=%s", username, roomID)
	// Step 1: 先获取聊天室对象
	chatroom, err := GetChatroom(roomID)
	if err != nil {
		log.Log.Warnf("聊天室不存在，加入失败: room_id=%s", roomID)
		return fmt.Errorf("聊天室不存在: %w", err)
	}

	// Step 2: 检查用户是否已经存在
	for _, u := range chatroom.Users {
		if u == username {
			log.Log.Infof("用户已在聊天室中: user=%s, room=%s", username, roomID)
			return nil
		}
	}
	log.Log.Infof("将用户添加到聊天室: user=%s, room=%s", username, roomID)
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
	if err != nil {
		log.Log.Errorf("写入用户更新失败: %v", err)
	} else {
		log.Log.Infof("用户加入聊天室成功: user=%s, room=%s", username, roomID)
	}
	return err
}

func RemoveUserFromChatroom(username, roomID string) error {
	log.Log.Infof("尝试将用户移出聊天室: user=%s, room=%s", username, roomID)
	room, err := GetChatroom(roomID)
	if err != nil {
		log.Log.Warnf("聊天室不存在，移除失败: room_id=%s", roomID)
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
	log.Log.Infof("将用户移出聊天室: user=%s, room=%s", username, roomID)
	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(ChatroomTableName),
		Item:      item,
	})
	if err != nil {
		log.Log.Errorf("移除用户失败: %v", err)
	} else {
		log.Log.Infof("用户移除成功: user=%s, room=%s", username, roomID)
	}
	return err
}

func GetChatroomsByUsername(username string) ([]Chatroom, error) {
	log.Log.Infof("查询用户加入的所有聊天室: user=%s", username)
	var results []Chatroom

	// 全表扫描
	output, err := DB.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(ChatroomTableName),
	})
	if err != nil {
		log.Log.Errorf("全表扫描失败: %v", err)
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
	log.Log.Infof("用户加入的聊天室总数: %d", len(results))
	return results, nil
}
