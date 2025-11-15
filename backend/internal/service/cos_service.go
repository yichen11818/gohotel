package service

import (
	"context"
	"fmt"
	"gohotel/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// CosService 腾讯云对象存储服务
type CosService struct {
	client *cos.Client
	config *config.COSConfig
}

// NewCosService 创建 COS 服务实例
func NewCosService(cfg *config.COSConfig) (*CosService, error) {
	if cfg.SecretID == "" || cfg.SecretKey == "" || cfg.BaseURL == "" {
		return nil, fmt.Errorf("COS 配置不完整，请检查环境变量")
	}

	// 解析 Bucket URL
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("COS BaseURL 格式错误: %v", err)
	}

	// 创建 COS 客户端
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &CosService{
		client: client,
		config: cfg,
	}, nil
}

// UploadFile 上传文件到 COS
// 参数:
//   - file: 上传的文件
//   - folder: 存储文件夹路径，如 "images", "documents"
//   - filename: 文件名（可选，如果为空则使用原始文件名）
//
// 返回:
//   - 文件在 COS 中的完整 URL
//   - 错误信息
func (s *CosService) UploadFile(file *multipart.FileHeader, folder string, filename string) (string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 如果没有指定文件名，使用原始文件名
	if filename == "" {
		filename = file.Filename
	}

	// 生成唯一的文件名（添加时间戳避免重名）
	timestamp := time.Now().Format("20060102150405")
	ext := path.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	uniqueFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)

	// 构建对象键（文件在 COS 中的路径）
	objectKey := path.Join(folder, uniqueFilename)

	// 上传文件到 COS
	_, err = s.client.Object.Put(context.Background(), objectKey, src, nil)
	if err != nil {
		return "", fmt.Errorf("上传文件到 COS 失败: %v", err)
	}

	// 返回文件的完整 URL
	fileURL := fmt.Sprintf("%s/%s", s.config.BaseURL, objectKey)
	return fileURL, nil
}

// UploadFileFromReader 从 io.Reader 上传文件到 COS
// 参数:
//   - reader: 文件内容的 Reader
//   - folder: 存储文件夹路径
//   - filename: 文件名
//
// 返回:
//   - 文件在 COS 中的完整 URL
//   - 错误信息
func (s *CosService) UploadFileFromReader(reader io.Reader, folder string, filename string) (string, error) {
	// 生成唯一的文件名
	timestamp := time.Now().Format("20060102150405")
	ext := path.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	uniqueFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)

	// 构建对象键
	objectKey := path.Join(folder, uniqueFilename)

	// 上传文件到 COS
	_, err := s.client.Object.Put(context.Background(), objectKey, reader, nil)
	if err != nil {
		return "", fmt.Errorf("上传文件到 COS 失败: %v", err)
	}

	// 返回文件的完整 URL
	fileURL := fmt.Sprintf("%s/%s", s.config.BaseURL, objectKey)
	return fileURL, nil
}

// DeleteFile 从 COS 删除文件
// 参数:
//   - fileURL: 文件的完整 URL
//
// 返回:
//   - 错误信息
func (s *CosService) DeleteFile(fileURL string) error {
	// 从 URL 中提取对象键
	objectKey, err := s.extractObjectKey(fileURL)
	if err != nil {
		return err
	}

	// 删除文件
	_, err = s.client.Object.Delete(context.Background(), objectKey)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// GetFileURL 获取文件的访问 URL
// 参数:
//   - objectKey: 对象键（文件在 COS 中的路径）
//
// 返回:
//   - 文件的完整 URL
func (s *CosService) GetFileURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", s.config.BaseURL, objectKey)
}

// GetPresignedURL 获取文件的预签名 URL（用于临时访问私有文件）
// 参数:
//   - objectKey: 对象键
//   - expireDuration: 过期时间
//
// 返回:
//   - 预签名 URL
//   - 错误信息
func (s *CosService) GetPresignedURL(objectKey string, expireDuration time.Duration) (string, error) {
	presignedURL, err := s.client.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet,
		objectKey,
		s.config.SecretID,
		s.config.SecretKey,
		expireDuration,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("生成预签名 URL 失败: %v", err)
	}

	return presignedURL.String(), nil
}

// extractObjectKey 从完整 URL 中提取对象键
func (s *CosService) extractObjectKey(fileURL string) (string, error) {
	// 解析 URL
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("URL 格式错误: %v", err)
	}

	// 提取路径部分（去掉开头的 "/"）
	objectKey := strings.TrimPrefix(u.Path, "/")
	if objectKey == "" {
		return "", fmt.Errorf("无法从 URL 中提取对象键")
	}

	return objectKey, nil
}

// FileExists 检查文件是否存在
// 参数:
//   - objectKey: 对象键
//
// 返回:
//   - 文件是否存在
//   - 错误信息
func (s *CosService) FileExists(objectKey string) (bool, error) {
	_, err := s.client.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		// 如果是 404 错误，说明文件不存在
		if cos.IsNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("检查文件是否存在失败: %v", err)
	}

	return true, nil
}
