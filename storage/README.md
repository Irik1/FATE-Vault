# Storage Service

A standalone storage service using MinIO for storing images and other files. This service operates independently and can continue to serve files even if other services are down.

## Features

- **File Upload**: Upload files (images, documents, etc.) via multipart form data
- **File Download**: Retrieve files by filename
- **File Deletion**: Delete files from storage
- **File Listing**: List all files in storage (with optional folder filtering)
- **Independent Operation**: Runs separately from backend and web services
- **MinIO Integration**: Uses MinIO for object storage

## Prerequisites

- Go 1.25.5 or later
- MinIO server (can be run via Docker Compose)

## Configuration

Copy `.env.example` to `.env` and configure:

```env
MINIO_ENDPOINT=localhost:9000
MINIO_USER=minioadmin
MINIO_PASSWORD=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=fate-vault
STORAGE_PORT=8081
```

## Running Locally

### Option 1: Using Docker Compose (Recommended)

This will start both MinIO and the storage service:

```bash
docker-compose up -d
```

**Note**: Make sure Docker Desktop is running before executing this command.

### Option 2: Using Startup Script

1. Start MinIO using the provided script:
```bash
./start-minio.sh
```

2. Wait a few seconds for MinIO to be ready, then run the storage service:
```bash
go run .
```

### Option 3: Manual Setup

1. **Start MinIO server** (make sure Docker is running):
```bash
docker run -d --name fate-vault-minio -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_PASSWORD=minioadmin \
  -v minio_data:/data \
  minio/minio server /data --console-address ":9001"
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Run the service**:
```bash
go run .
```

## Troubleshooting

### Connection Refused Error

If you see `connection refused` when starting the storage service:

1. **Check if Docker is running**:
   ```bash
   docker info
   ```
   If this fails, start Docker Desktop.

2. **Check if MinIO is running**:
   ```bash
   docker ps | grep minio
   ```
   If MinIO is not running, start it using one of the methods above.

3. **Wait for MinIO to be ready**: After starting MinIO, wait 5-10 seconds before running the storage service. The service will automatically retry connecting up to 5 times.

4. **Verify MinIO is accessible**:
   - API: http://localhost:9000
   - Console: http://localhost:9001 (login with minioadmin/minioadmin)

## API Endpoints

> **📖 For detailed curl examples and integration code, see [API_EXAMPLES.md](./API_EXAMPLES.md)**

### Health Check
```
GET /health
```
Returns service health status.

### Upload File
```
POST /upload
Content-Type: multipart/form-data

Form fields:
- file: The file to upload (required)
- folder: Optional folder/path prefix (optional)
```

Response:
```json
{
  "message": "File uploaded successfully",
  "filename": "images/character.jpg",
  "size": 12345,
  "contentType": "image/jpeg",
  "url": "/download/images/character.jpg"
}
```

### Download File
```
GET /download/:filename
```

Returns the file content with appropriate headers.

### Delete File
```
DELETE /delete/:filename
```

Response:
```json
{
  "message": "File deleted successfully",
  "filename": "images/character.jpg"
}
```

### List Files
```
GET /list?folder=images
```

Query parameters:
- `folder`: Optional folder prefix to filter files

Response:
```json
{
  "files": [
    {
      "name": "images/character.jpg",
      "size": 12345,
      "lastModified": "2024-01-01T00:00:00Z",
      "contentType": "image/jpeg",
      "url": "/download/images/character.jpg"
    }
  ],
  "count": 1
}
```

## Integration with Other Services

Other services can interact with the storage service by making HTTP requests:

### Example: Upload from Backend Service

```go
func uploadToStorage(filePath string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", filepath.Base(filePath))
    if err != nil {
        return err
    }
    io.Copy(part, file)
    writer.Close()

    req, err := http.NewRequest("POST", "http://localhost:8081/upload", body)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}
```

### Example: Get File URL

```go
func getFileURL(filename string) string {
    return fmt.Sprintf("http://localhost:8081/download/%s", filename)
}
```

## MinIO Console

Access the MinIO web console at `http://localhost:9001` to:
- Browse files
- Manage buckets
- Configure access policies
- Monitor storage usage

Default credentials:
- Username: `minioadmin`
- Password: `minioadmin`

## Production Considerations

1. **Security**: Change default MinIO credentials
2. **SSL/TLS**: Enable SSL for MinIO in production
3. **Access Control**: Configure MinIO bucket policies
4. **Backup**: Set up regular backups of MinIO data
5. **Monitoring**: Add logging and monitoring
6. **Rate Limiting**: Consider adding rate limiting to API endpoints

