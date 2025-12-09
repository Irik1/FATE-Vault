package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKeyID := os.Getenv("MINIO_USER")
	secretAccessKey := os.Getenv("MINIO_PASSWORD")
	useSSL := os.Getenv("MINIO_USE_SSL")
	secure := useSSL == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	minioClient = client

	// Test connection with retry logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Attempting to connect to MinIO at %s...\n", endpoint)

	maxRetries := 5
	retryDelay := 2 * time.Second
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = client.ListBuckets(ctx)
		cancel()

		if err == nil {
			fmt.Println("Connected to MinIO successfully")
			break
		}

		lastErr = err
		if i < maxRetries-1 {
			fmt.Printf("Connection attempt %d/%d failed: %v. Retrying in %v...\n", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if lastErr != nil {
		log.Fatalf("Failed to connect to MinIO after %d attempts: %v\n\n"+
			"Please ensure MinIO is running. You can start it with:\n"+
			"  1. Docker Compose: cd storage && docker-compose up -d\n"+
			"  2. Docker: docker run -d -p 9000:9000 -p 9001:9001 -e MINIO_ROOT_USER=minioadmin -e MINIO_PASSWORD=minioadmin minio/minio server /data --console-address \":9001\"\n"+
			"  3. Or set MINIO_ENDPOINT to point to your MinIO instance", maxRetries, lastErr)
	}

	// Ensure default bucket exists
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	err = ensureBucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket '%s' is ready\n", bucketName)

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "storage",
		})
	})

	// Storage endpoints
	router.POST("/upload", uploadFile)
	router.GET("/download/:filename", downloadFile)
	router.DELETE("/delete/:filename", deleteFile)
	router.GET("/list", listFiles)

	port := os.Getenv("STORAGE_PORT")

	fmt.Printf("Storage service starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func ensureBucketExists(ctx context.Context, bucketName string) error {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Created bucket: %s\n", bucketName)
	}

	return nil
}
