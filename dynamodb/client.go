package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DB *dynamodb.Client

func InitDB() {
	env := os.Getenv("DYNAMODB_ENV")
	if env == "" {
		env = "local" // 默认环境
	}

	region := "us-west-2" // 可以放进 env 里也行
	var cfg aws.Config
	var err error

	if env == "local" {
		log.Println("🌱 连接本地 DynamoDB (local mode)")

		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					URL:           "http://localhost:8000", // DynamoDB Local 地址
					SigningRegion: region,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
			// Add dummy credentials for local mode
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		)

		if err != nil {
			log.Fatal("❌ 加载本地 DynamoDB 配置失败:", err)
		}

	} else if env == "aws" {
		log.Println("🚀 连接 AWS DynamoDB（真实云服务）")

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			log.Fatal("❌ 加载 AWS 配置失败:", err)
		}
	} else {
		log.Fatalf("❌ 未知 DYNAMODB_ENV：%s", env)
	}

	DB = dynamodb.NewFromConfig(cfg)
	log.Println("✅ DynamoDB 客户端初始化成功")
}

func CreateAllTables() {
	//CreateUserTable()
	//CreateChatroomTable()
	CreateMessageTable()
}
