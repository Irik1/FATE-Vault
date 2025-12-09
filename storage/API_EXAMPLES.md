# Storage Service API Examples

This document provides curl examples for all endpoints in the storage service.

**Base URL**: `http://localhost:8081` (default port)

---

## 1. Health Check

Check if the storage service is running.

```bash
curl -X GET http://localhost:8081/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "storage"
}
```

---

## 2. Upload File

Upload a file to storage. Files can be organized in folders.

### Basic Upload

```bash
curl -X POST http://localhost:8081/upload \
  -F "file=@/path/to/your/image.jpg"
```

### Upload to Specific Folder

```bash
curl -X POST http://localhost:8081/upload \
  -F "file=@/path/to/your/image.jpg" \
  -F "folder=images/characters"
```

### Upload File with Spaces and Special Characters

For files with spaces or special characters (like Cyrillic), use quotes around the path:

```bash
# Upload PDF from Downloads folder (Mac/Linux)
curl -X POST http://localhost:8081/upload \
  -F "file=@$HOME/Downloads/FATE Warcraft прегены.pdf" \
  -F "folder=documents"

# Or with full path
curl -X POST http://localhost:8081/upload \
  -F "file=@/Users/yourusername/Downloads/FATE Warcraft прегены.pdf" \
  -F "folder=documents"

# Upload to root (no folder)
curl -X POST http://localhost:8081/upload \
  -F "file=@$HOME/Downloads/FATE Warcraft прегены.pdf"
```

### Upload Multiple Files (call multiple times)

```bash
# Upload first file
curl -X POST http://localhost:8081/upload \
  -F "file=@character1.jpg" \
  -F "folder=images/characters"

# Upload second file
curl -X POST http://localhost:8081/upload \
  -F "file=@character2.png" \
  -F "folder=images/characters"
```

### Upload with Verbose Output

```bash
curl -v -X POST http://localhost:8081/upload \
  -F "file=@document.pdf" \
  -F "folder=documents"
```

**Response:**
```json
{
  "message": "File uploaded successfully",
  "filename": "images/characters/image.jpg",
  "size": 123456,
  "contentType": "image/jpeg",
  "url": "/download/images/characters/image.jpg"
}
```

---

## 3. Download File

Download a file from storage.

### Basic Download

```bash
curl -X GET http://localhost:8081/download/image.jpg \
  -o downloaded_image.jpg
```

### Download from Folder

```bash
curl -X GET http://localhost:8081/download/images/characters/character1.jpg \
  -o character1.jpg
```

### Download and Display in Terminal (for text files)

```bash
curl -X GET http://localhost:8081/download/documents/readme.txt
```

### Download with Progress Bar

```bash
curl -X GET http://localhost:8081/download/large-file.zip \
  --progress-bar \
  -o large-file.zip
```

### Download and Save with Original Name

```bash
# Get filename from response headers
curl -X GET http://localhost:8081/download/images/photo.jpg \
  -O -J
```

---

## 4. Delete File

Delete a file from storage.

### Delete File in Root

```bash
curl -X DELETE http://localhost:8081/delete/image.jpg
```

### Delete File in Folder

```bash
curl -X DELETE http://localhost:8081/delete/images/characters/character1.jpg
```

**Response:**
```json
{
  "message": "File deleted successfully",
  "filename": "images/characters/character1.jpg"
}
```

---

## 5. List Files

List all files in storage, optionally filtered by folder.

### List All Files

```bash
curl -X GET http://localhost:8081/list
```

### List Files in Specific Folder

```bash
curl -X GET "http://localhost:8081/list?folder=images/characters"
```

### Pretty Print JSON Response

```bash
curl -X GET http://localhost:8081/list | jq
```

**Response:**
```json
{
  "files": [
    {
      "name": "images/characters/character1.jpg",
      "size": 123456,
      "lastModified": "2024-12-08T16:30:00Z",
      "contentType": "image/jpeg",
      "url": "/download/images/characters/character1.jpg"
    },
    {
      "name": "images/characters/character2.png",
      "size": 234567,
      "lastModified": "2024-12-08T16:31:00Z",
      "contentType": "image/png",
      "url": "/download/images/characters/character2.png"
    }
  ],
  "count": 2
}
```

---

## Complete Workflow Example

Here's a complete example of uploading, listing, downloading, and deleting a file:

```bash
# 1. Upload a file
curl -X POST http://localhost:8081/upload \
  -F "file=@test-image.jpg" \
  -F "folder=test"

# Response will include the filename, e.g., "test/test-image.jpg"

# 2. List all files in test folder
curl -X GET "http://localhost:8081/list?folder=test"

# 3. Download the file
curl -X GET http://localhost:8081/download/test/test-image.jpg \
  -o downloaded-test-image.jpg

# 4. Delete the file
curl -X DELETE http://localhost:8081/delete/test/test-image.jpg

# 5. Verify deletion (should return empty list)
curl -X GET "http://localhost:8081/list?folder=test"
```

---

## Error Handling Examples

### File Not Found (404)

```bash
curl -X GET http://localhost:8081/download/nonexistent.jpg
```

**Response:**
```json
{
  "error": "File not found: ..."
}
```

### Missing File in Upload (400)

```bash
curl -X POST http://localhost:8081/upload
```

**Response:**
```json
{
  "error": "No file provided: ..."
}
```

### Invalid Endpoint (404)

```bash
curl -X GET http://localhost:8081/invalid-endpoint
```

---

## Advanced Examples

### Upload with Custom Headers

```bash
curl -X POST http://localhost:8081/upload \
  -H "Authorization: Bearer your-token" \
  -F "file=@image.jpg" \
  -F "folder=images"
```

### Download with Range Request (for partial downloads)

```bash
curl -X GET http://localhost:8081/download/large-file.zip \
  -H "Range: bytes=0-1023" \
  -o partial-file.zip
```

### Check File Exists (using HEAD request - if implemented)

```bash
curl -I http://localhost:8081/download/image.jpg
```

---

## Integration Examples

### From Go Backend

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

func uploadFile(filePath string, folder string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    
    part, err := writer.CreateFormFile("file", file.Name())
    if err != nil {
        return err
    }
    io.Copy(part, file)
    
    if folder != "" {
        writer.WriteField("folder", folder)
    }
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

### From JavaScript/TypeScript

```javascript
async function uploadFile(file, folder = '') {
  const formData = new FormData();
  formData.append('file', file);
  if (folder) {
    formData.append('folder', folder);
  }

  const response = await fetch('http://localhost:8081/upload', {
    method: 'POST',
    body: formData
  });

  return await response.json();
}

// Usage
const fileInput = document.querySelector('input[type="file"]');
const file = fileInput.files[0];
uploadFile(file, 'images/characters')
  .then(result => console.log('Uploaded:', result))
  .catch(error => console.error('Error:', error));
```

### From Python

```python
import requests

def upload_file(file_path, folder=''):
    url = 'http://localhost:8081/upload'
    files = {'file': open(file_path, 'rb')}
    data = {'folder': folder} if folder else {}
    
    response = requests.post(url, files=files, data=data)
    return response.json()

# Usage
result = upload_file('image.jpg', 'images/characters')
print(result)
```

---

## Testing Script

Save this as `test-storage-api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8081"

echo "1. Health Check..."
curl -s "$BASE_URL/health" | jq
echo ""

echo "2. Upload test file..."
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/upload" \
  -F "file=@test.txt" \
  -F "folder=test")
echo "$UPLOAD_RESPONSE" | jq
FILENAME=$(echo "$UPLOAD_RESPONSE" | jq -r '.filename')
echo ""

echo "3. List files..."
curl -s "$BASE_URL/list?folder=test" | jq
echo ""

echo "4. Download file..."
curl -s -X GET "$BASE_URL/download/$FILENAME" -o downloaded_test.txt
echo "Downloaded to downloaded_test.txt"
echo ""

echo "5. Delete file..."
curl -s -X DELETE "$BASE_URL/delete/$FILENAME" | jq
echo ""

echo "6. Verify deletion..."
curl -s "$BASE_URL/list?folder=test" | jq
```

Make it executable and run:
```bash
chmod +x test-storage-api.sh
./test-storage-api.sh
```

