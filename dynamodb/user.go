package dynamodb

import (
	log "chatroom-api/logger"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Username string `dynamodbav:"username"` // 👈 主键
	Password string `dynamodbav:"password"`
}

var UserTableName = "users"

func CreateUser(user User) error {
	log.Log.Infof("尝试创建用户: username=%s", user.Username)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Log.Errorf("用户数据序列化失败: %v", err)
		return err
	}

	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           &UserTableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"), // 防止重复注册
	})
	if err != nil {
		log.Log.Warnf("用户创建失败: username=%s, err=%v", user.Username, err)
	} else {
		log.Log.Infof("用户创建成功: username=%s", user.Username)
	}
	return err
}

func GetUserByUsername(username string) (*User, error) {
	log.Log.Infof("尝试获取用户: username=%s", username)
	out, err := DB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &UserTableName,
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		log.Log.Errorf("查询用户失败: username=%s, err=%v", username, err)
		return nil, errors.New("user not found")
	}
	if out.Item == nil {
		log.Log.Warnf("用户不存在: username=%s", username)
		return nil, errors.New("user not found")
	}

	var user User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		log.Log.Errorf("用户反序列化失败: username=%s, err=%v", username, err)
		return nil, err
	}

	log.Log.Infof("成功获取用户信息: username=%s", user.Username)
	return &user, nil
}

func CreateUserTable() error {
	log.Log.Info("开始创建 users 表")
	_, err := DB.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(UserTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("username"),
				AttributeType: types.ScalarAttributeTypeS, // String 类型
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("username"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest, // 免费账号推荐按需计费
	})
	if err != nil {
		// 容错：如果表已存在，不返回错误
		var rne *types.ResourceInUseException
		if errors.As(err, &rne) {
			log.Log.Info("⚠️ 用户表 [%s] 已存在，跳过创建", UserTableName)
			return nil
		}

		// 其余错误要向上传递
		return fmt.Errorf("创建用户表 [%s] 失败: %w", UserTableName, err)
	}
	log.Log.Info("users 表创建成功")
	return nil
}
