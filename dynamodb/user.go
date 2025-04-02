package dynamodb

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type User struct {
	Username string `dynamodbav:"username"` // 👈 主键
	Password string `dynamodbav:"password"`
}

var UserTableName = "users"

func CreateUser(user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = DB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           &UserTableName,
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"), // 防止重复注册
	})
	return err
}

func GetUserByUsername(username string) (*User, error) {
	out, err := DB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &UserTableName,
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil || out.Item == nil {
		return nil, errors.New("user not found")
	}

	var user User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUserTable() {
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
		log.Fatalf("❌ 创建 users 表失败: %v", err)
	}
	log.Println("✅ 用户表创建成功")
}
