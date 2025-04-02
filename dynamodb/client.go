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
	endpoint := os.Getenv("DYNAMODB_ENDPOINT") // 本地模式會設這個
	region := os.Getenv("DYNAMODB_REGION")
	if region == "" {
		region = "us-west-2" // fallback
	}
	var cfg aws.Config
	var err error

	if endpoint != "" {
		log.Println("🌱 连接本地 DynamoDB (local mode)")

		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{
					URL:           endpoint, // DynamoDB Local 地址
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

	} else {
		log.Println("🚀 连接 AWS DynamoDB（真实云服务）")

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			log.Fatal("❌ 加载 AWS 配置失败:", err)
		}
	}

	DB = dynamodb.NewFromConfig(cfg)
	log.Println("✅ DynamoDB 客户端初始化成功")
}

func CreateAllTables() {
	//CreateUserTable()
	//CreateChatroomTable()
	CreateMessageTable()
}
