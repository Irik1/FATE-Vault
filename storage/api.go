package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func uploadFile(c *gin.Context) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided: " + err.Error()})
		return
	}
	defer file.Close()

	folder := c.PostForm("folder")
	if folder != "" {
		folder = strings.Trim(folder, "/") + "/"
	}

	filename := fmt.Sprintf("file_%d", time.Now().Unix())
	objectName := folder + filename

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = minioClient.PutObject(ctx, bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Return file info
	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"filename":    objectName,
		"size":        header.Size,
		"contentType": contentType,
		"url":         fmt.Sprintf("/download/%s", objectName),
	})
}

func downloadFile(c *gin.Context) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	object, err := minioClient.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found: " + err.Error()})
		return
	}
	defer object.Close()

	objInfo, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info: " + err.Error()})
		return
	}

	// Set headers
	c.Header("Content-Type", objInfo.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(filename)))
	c.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// Stream file to response
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream file: " + err.Error()})
		return
	}
}

func deleteFile(c *gin.Context) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := minioClient.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File deleted successfully",
		"filename": filename,
	})
}

func listFiles(c *gin.Context) {
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "fate-vault"
	}

	// Get optional folder prefix
	folder := c.Query("folder")
	if folder != "" {
		folder = strings.Trim(folder, "/") + "/"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    folder,
		Recursive: true,
	})

	var files []map[string]interface{}
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files: " + object.Err.Error()})
			return
		}

		files = append(files, map[string]interface{}{
			"name":         object.Key,
			"size":         object.Size,
			"lastModified": object.LastModified,
			"contentType":  object.ContentType,
			"url":          fmt.Sprintf("/download/%s", object.Key),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"count": len(files),
	})
}
