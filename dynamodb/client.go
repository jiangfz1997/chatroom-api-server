package dynamodb

import (
	log "chatroom-api/logger"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		region = "us-west-2"
		log.Log.Warn("⚠️ 未设置 DYNAMODB_REGION，默认使用 us-west-2")
	} else {
		log.Log.Infof("✅ 检测到 DYNAMODB_REGION: %s", region)
	}
	var cfg aws.Config
	var err error

	if endpoint != "" {
		log.Log.Info("连接本地 DynamoDB (local mode)")
		log.Log.Infof("当前 DynamoDB Endpoint: %s", endpoint)
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				log.Log.Infof("📍 自定义解析器：服务 [%s]，区域 [%s]", service, region)
				return aws.Endpoint{
					URL:           endpoint, // DynamoDB Local 地址
					SigningRegion: region,
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested %s", service)
		})

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
			// Add dummy credentials for local mode
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		)

		if err != nil {
			log.Log.Fatalf("加载本地 DynamoDB 配置失败: %v", err)
		}
		log.Log.Info("本地 DynamoDB 配置加载成功")
	} else {
		log.Log.Info("连接 AWS DynamoDB（真实云服务）")

		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			log.Log.Fatalf("加载 AWS 配置失败: %v", err)
		}
		log.Log.Info("AWS配置加载成功")
	}

	DB = dynamodb.NewFromConfig(cfg)
	log.Log.Info("DynamoDB 客户端初始化成功")
}

func CreateAllTables() error {
	var errs []error

	if err := CreateUserTable(); err != nil {
		errs = append(errs, fmt.Errorf("CreateUserTable failed: %w", err))
	}
	if err := CreateChatroomTable(); err != nil {
		errs = append(errs, fmt.Errorf("CreateChatroomTable failed: %w", err))
	}
	if err := CreateMessageTable(); err != nil {
		errs = append(errs, fmt.Errorf("CreateMessageTable failed: %w", err))
	}
	if len(errs) > 0 {
		errMsg := "❌ Table creation encountered errors:\n"
		for _, e := range errs {
			errMsg += " - " + e.Error() + "\n"
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}
