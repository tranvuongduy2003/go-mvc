# File Upload System Setup Guide

## Tổng quan

Hệ thống upload file tổng quát sử dụng MinIO làm object storage và Traefik làm CDN proxy, hỗ trợ nhiều loại file khác nhau như avatar, documents, images, attachments với luồng upload được thiết kế để lưu trữ file an toàn và phục vụ nhanh chóng.

## Kiến trúc

```
User Upload → API (api.localhost) → MinIO Storage → CDN (cdn.localhost)
```

### Thành phần:
- **MinIO**: S3-compatible object storage
- **Traefik**: Reverse proxy làm CDN cho static files
- **API Backend**: Xử lý upload logic và database
- **File Types Support**: Avatars, Documents, Images, Attachments, Media files

## Các loại file được hỗ trợ

### 1. Avatar Upload (User Profile Images)
- **Endpoint**: `POST /api/v1/users/:id/avatar`
- **Supported formats**: PNG, JPEG, JPG, GIF
- **Size limit**: 5MB
- **Validation**: Image format và file size

### 2. Document Upload (Extensible)
- **Pattern**: `POST /api/v1/{entity}/{id}/{file-type}`
- **Examples**:
  - `POST /api/v1/users/:id/documents`
  - `POST /api/v1/projects/:id/attachments`
  - `POST /api/v1/posts/:id/images`
- **Supported formats**: PDF, DOC, DOCX, TXT, etc.

### 3. Media Upload
- **Audio**: MP3, WAV, AAC
- **Video**: MP4, AVI, MOV
- **Archives**: ZIP, RAR, 7Z

## Docker Compose Setup

### 1. Services được thêm vào docker-compose.yml:

```yaml
# Traefik Reverse Proxy
traefik:
  image: traefik:v3.0
  container_name: dev-traefik
  restart: unless-stopped
  ports:
    - "80:80"      # HTTP
    - "8081:8080"  # Traefik Dashboard
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock:ro
    - ./configs/traefik:/etc/traefik/dynamic:ro

# MinIO Object Storage  
minio:
  image: minio/minio:latest
  container_name: dev-minio
  environment:
    MINIO_ROOT_USER: minioadmin
    MINIO_ROOT_PASSWORD: minioadmin
  ports:
    - "9000:9000"  # API
    - "9001:9001"  # Console
  volumes:
    - minio_data:/data
  labels:
    - traefik.http.routers.minio-api.rule=Host(`cdn.localhost`)
    - traefik.http.routers.minio-console.rule=Host(`minio.localhost`)
```

### 2. Traefik Configuration:

```yaml
# configs/traefik/dynamic.yml
http:

## Configuration

### File Storage Config (`configs/development.yaml`):

```yaml
external:
  file_storage:
    endpoint: "localhost:9000"
    access_key_id: "minioadmin"
    secret_access_key: "minioadmin"
    bucket_name: "uploads"
    cdn_url: "http://cdn.localhost"
    use_ssl: false
```

## Database Schema

### File storage pattern cho các entities:

```sql
-- Users table (avatar example)
ALTER TABLE users 
ADD COLUMN avatar_file_key VARCHAR(500),
ADD COLUMN avatar_cdn_url VARCHAR(1000);

-- Generic file attachment pattern
CREATE TABLE file_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,  -- 'users', 'projects', 'posts', etc.
    entity_id UUID NOT NULL,           -- ID of the parent entity
    file_type VARCHAR(50) NOT NULL,    -- 'avatar', 'document', 'image', etc.
    file_key VARCHAR(500) NOT NULL,    -- S3 object key
    file_name VARCHAR(255) NOT NULL,   -- Original filename
    file_size BIGINT NOT NULL,         -- File size in bytes
    mime_type VARCHAR(100) NOT NULL,   -- MIME type
    cdn_url VARCHAR(1000) NOT NULL,    -- CDN access URL
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for efficient queries
CREATE INDEX idx_file_attachments_entity ON file_attachments(entity_type, entity_id, file_type);
```

## API Usage

### 1. Upload Avatar (Specific Implementation):

```http
POST /api/v1/users/{id}/avatar
Content-Type: multipart/form-data

Form Data:
- avatar: [file] (max 5MB, image formats only)
```

### 2. Generic File Upload Pattern:

```http
POST /api/v1/{entity}/{id}/files
Content-Type: multipart/form-data

Form Data:
- file: [file]
- file_type: [string] (e.g., "document", "image", "attachment")
- description: [string] (optional)
```

### 3. Examples of Extended Usage:

```bash
# Upload project document
POST /api/v1/projects/{project_id}/files
Form: file=document.pdf, file_type="document"

# Upload post images
POST /api/v1/posts/{post_id}/files  
Form: file=image.jpg, file_type="image"

# Upload user profile documents
POST /api/v1/users/{user_id}/files
Form: file=resume.pdf, file_type="resume"
```

### Response:

```json
{
  "success": true,
  "data": {
    "id": "user-uuid",
    "email": "user@example.com", 
    "name": "User Name",
    "phone": "+1234567890",
    "avatar_url": "http://cdn.localhost/uploads/avatars/user-uuid/file.jpg",
    "is_active": true,
    "created_at": "2025-09-27T00:00:00Z",
    "updated_at": "2025-09-27T00:00:00Z"
  }
}
```

## Cách sử dụng

### 1. Khởi động services:

```bash
# Start all services including MinIO and Traefik
docker-compose up -d

# Verify services are running
docker-compose ps
```

### 2. Setup hosts (MacOS/Linux):

```bash
# Add to /etc/hosts
echo "127.0.0.1 api.localhost" | sudo tee -a /etc/hosts
echo "127.0.0.1 cdn.localhost" | sudo tee -a /etc/hosts  
echo "127.0.0.1 minio.localhost" | sudo tee -a /etc/hosts
```

### 3. Chạy migration:

```bash
# Run database migrations to add avatar columns
make migrate-up
# or
go run ./cmd/migrate/main.go -action=up
```

### 4. Start API server:

```bash
# Start the Go API server
make run
# or
go run ./cmd/main.go
```

### 5. Access points:

- **API**: http://api.localhost/api/v1
- **CDN/Files**: http://cdn.localhost/uploads/...
- **MinIO Console**: http://minio.localhost (minioadmin/minioadmin)
- **Traefik Dashboard**: http://localhost:8081

## Implementation Pattern

### 1. File Upload Command Pattern:

```go
// Example: Document upload command
type UploadDocumentCommand struct {
    EntityType string        // "projects", "users", etc.
    EntityID   string        // Entity UUID
    FileType   string        // "document", "image", etc.
    File       *multipart.FileHeader
}

type UploadDocumentCommandHandler struct {
    fileStorage ports.FileStorageService
    repository  ports.FileAttachmentRepository
}

func (h *UploadDocumentCommandHandler) Handle(cmd UploadDocumentCommand) (*FileAttachment, error) {
    // 1. Validate file
    if err := h.validateFile(cmd.File, cmd.FileType); err != nil {
        return nil, err
    }
    
    // 2. Upload to storage
    fileKey := fmt.Sprintf("%s/%s/%s", cmd.EntityType, cmd.EntityID, generateFileName(cmd.File))
    cdnURL, err := h.fileStorage.UploadFile(fileKey, cmd.File)
    if err != nil {
        return nil, err
    }
    
    // 3. Save to database
    attachment := &FileAttachment{
        EntityType: cmd.EntityType,
        EntityID:   cmd.EntityID,
        FileType:   cmd.FileType,
        FileKey:    fileKey,
        CDNURL:     cdnURL,
        // ... other fields
    }
    
    return h.repository.Create(attachment)
}
```

### 2. Generic File Handler:

```go
func (h *FileHandler) UploadFile(c *gin.Context) {
    entityType := c.Param("entity")
    entityID := c.Param("id") 
    fileType := c.PostForm("file_type")
    
    file, err := c.FormFile("file")
    if err != nil {
        response.Error(c, http.StatusBadRequest, "No file provided")
        return
    }
    
    cmd := UploadFileCommand{
        EntityType: entityType,
        EntityID:   entityID,
        FileType:   fileType,
        File:       file,
    }
    
    result, err := h.uploadHandler.Handle(cmd)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    response.Success(c, result)
}
```

## Testing File Upload

### 1. Avatar Upload:

```bash
# Upload avatar for user
curl -X POST \
  http://api.localhost/api/v1/users/{user-id}/avatar \
  -F "avatar=@/path/to/image.jpg" \
  -H "Authorization: Bearer {jwt-token}"
```

### 2. Generic File Upload Examples:

```bash
# Upload project document
curl -X POST \
  http://api.localhost/api/v1/projects/{project-id}/files \
  -F "file=@/path/to/document.pdf" \
  -F "file_type=document" \
  -H "Authorization: Bearer {jwt-token}"

# Upload user resume
curl -X POST \
  http://api.localhost/api/v1/users/{user-id}/files \
  -F "file=@/path/to/resume.pdf" \
  -F "file_type=resume" \
  -H "Authorization: Bearer {jwt-token}"

# Upload post images
curl -X POST \
  http://api.localhost/api/v1/posts/{post-id}/files \
  -F "file=@/path/to/image.jpg" \
  -F "file_type=image" \
  -H "Authorization: Bearer {jwt-token}"
```

### 3. Using Postman:

1. Method: POST
2. URL: `http://api.localhost/api/v1/{entity}/{id}/{endpoint}`
3. Headers: `Authorization: Bearer {jwt-token}`
4. Body: form-data with appropriate file fields

## File Structure

```
internal/
├── adapters/
│   └── external/
│       └── file_storage_service.go    # MinIO S3 client (supports all file types)
├── application/
│   ├── commands/
│   │   ├── user/
│   │   │   └── upload_avatar_command.go   # Avatar-specific command
│   │   └── shared/
│   │       └── upload_file_command.go     # Generic file upload command
│   └── dto/
│       ├── user/user_dto.go              # User with avatar_url
│       └── shared/file_attachment_dto.go  # Generic file attachment DTO
├── core/domain/
│   ├── user/user.go                      # User aggregate with Avatar
│   └── shared/
│       └── file_attachment.go           # Generic file attachment entity
└── di/
    ├── infrastructure.go                # File storage service provider
    └── modules/
        ├── user.go                      # User-specific handlers
        └── shared.go                    # Generic file upload handlers
```

## Monitoring & Troubleshooting

### Logs:

```bash
# Application logs
docker-compose logs app

# MinIO logs  
docker-compose logs minio

# Traefik logs
docker-compose logs traefik
```

### MinIO Console:

- URL: http://minio.localhost
- Credentials: minioadmin/minioadmin
- Browse uploaded files in `uploads` bucket

### Common Issues:

1. **File not accessible via CDN**: Check Traefik routing configuration
2. **Upload fails**: Verify MinIO credentials and bucket policy
3. **Large file upload**: Adjust nginx/traefik upload limits
4. **CORS issues**: Update Traefik middleware configuration

## Security Considerations

### File Validation:
- **File size limits**: Configured per file type (avatars: 5MB, documents: 20MB, etc.)
- **MIME type validation**: Only allow specific file types per category
- **File extension verification**: Double-check file extensions match MIME types
- **Malware scanning**: Consider integration with antivirus services for production

### Access Control:
- **Authentication required**: All upload endpoints require valid JWT tokens
- **Authorization checks**: Verify user permissions for specific entity uploads
- **Rate limiting**: Implement upload rate limits to prevent abuse
- **CDN security**: Traefik configuration includes security headers

### Production Deployment:
- **SSL/TLS**: Enable HTTPS for all file operations
- **CDN optimization**: Configure proper caching headers and compression
- **Backup strategy**: Implement MinIO backup and disaster recovery
- **Monitoring**: Track upload metrics, failed uploads, and storage usage

## Extending the System

### Adding New File Types:

1. **Create entity-specific commands**:
```go
// internal/application/commands/project/upload_document_command.go
type UploadDocumentCommand struct {
    ProjectID   uuid.UUID
    File        multipart.File
    FileHeader  *multipart.FileHeader
    FileType    string
    UploadedBy  uuid.UUID
}
```

2. **Implement command handlers**:
```go
type UploadDocumentCommandHandler struct {
    fileStorageService ports.FileStorageService
    projectRepository  ports.ProjectRepository
    logger             *zap.Logger
}

func (h *UploadDocumentCommandHandler) Handle(cmd *UploadDocumentCommand) (*dto.FileAttachmentDTO, error) {
    // Validate project ownership
    // Upload to MinIO with "documents" bucket
    // Update project with document reference
    // Return file metadata
}
```

3. **Create HTTP endpoints**:
```go
// internal/handlers/http/rest/project_handler.go
func (h *ProjectHandler) UploadDocument(c *gin.Context) {
    // Parse multipart form
    // Create upload command
    // Execute command handler
    // Return file information
}
```

### Database Schema Extensions:

```sql
-- Generic file attachments table
CREATE TABLE file_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(50) NOT NULL,     -- 'user', 'project', 'post', etc.
    entity_id UUID NOT NULL,
    file_type VARCHAR(50) NOT NULL,       -- 'avatar', 'document', 'image', etc.
    original_name VARCHAR(255) NOT NULL,
    file_key VARCHAR(500) NOT NULL,       -- MinIO object key
    cdn_url VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(entity_type, entity_id, file_type),  -- One file per type per entity
    FOREIGN KEY (uploaded_by) REFERENCES users(id)
);

-- Index for efficient queries
CREATE INDEX idx_file_attachments_entity ON file_attachments(entity_type, entity_id);
CREATE INDEX idx_file_attachments_type ON file_attachments(file_type);
```

This comprehensive file upload system provides a robust foundation that can be extended for any file upload needs in your Go MVC application.