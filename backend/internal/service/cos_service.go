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
	"strconv"
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
	return s.uploadFile(file, folder, filename, false)
}

// UploadTempFile 上传临时文件到 COS
// 参数:
//   - file: 上传的文件
//   - resourceType: 资源类型，如 "banner", "room", "user"
//
// 返回:
//   - 临时文件的完整 URL
//   - 错误信息
func (s *CosService) UploadTempFile(file *multipart.FileHeader, resourceType string) (string, error) {
	// 恢复临时文件夹路径，包含日期目录
	tempFolder := path.Join("tmp", resourceType, time.Now().Format("20060102"))
	return s.uploadFile(file, tempFolder, "", true)
}

// uploadFile 内部上传文件方法
func (s *CosService) uploadFile(file *multipart.FileHeader, folder string, filename string, isTemp bool) (string, error) {
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

	// 生成唯一的文件名（添加时间戳避免重名，移除随机字符串）
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

	// 返回文件的完整 URL，确保只有一个斜杠
	baseURL := strings.TrimRight(s.config.BaseURL, "/")
	fileURL := fmt.Sprintf("%s/%s", baseURL, objectKey)

	// 如果是临时文件，添加临时标识和过期时间
	if isTemp {
		expireTime := time.Now().Add(1 * time.Hour) // 1小时后过期
		fileURL = fmt.Sprintf("%s?temp=true&expire=%d", fileURL, expireTime.Unix())
	}

	return fileURL, nil
}

// ConfirmUpload 确认上传，将临时文件转为正式文件
// 参数:
//   - tempURL: 临时文件URL
//
// 返回:
//   - 正式文件的完整 URL
//   - 错误信息
func (s *CosService) ConfirmUpload(tempURL string) (string, error) {
	// 解析临时URL
	u, err := url.Parse(tempURL)
	if err != nil {
		return "", fmt.Errorf("URL 格式错误: %v", err)
	}

	// 检查是否为临时文件
	if u.Query().Get("temp") != "true" {
		return "", fmt.Errorf("不是临时文件URL")
	}

	// 检查是否过期
	expireStr := u.Query().Get("expire")
	expire, err := strconv.ParseInt(expireStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("无效的过期时间")
	}
	if time.Now().Unix() > expire {
		return "", fmt.Errorf("临时文件已过期")
	}

	// 获取临时对象键，确保路径正确（去除开头的/）
	tempObjectKey := strings.TrimPrefix(u.Path, "/")

	// 检查临时文件是否存在
	exists, err := s.FileExists(tempObjectKey)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("临时文件不存在")
	}

	// 解析临时文件路径，获取资源类型和文件名
	// 临时路径格式: tmp/resourceType/20060102/filename
	pathParts := strings.Split(tempObjectKey, "/")
	if len(pathParts) < 4 {
		return "", fmt.Errorf("无效的临时文件路径格式")
	}

	resourceType := pathParts[1]
	filename := pathParts[3] // 跳过日期目录，直接获取文件名

	// 构建正式对象键，使用资源类型的复数形式
	// 如: banner -> banners
	formalFolder := fmt.Sprintf("%ss", resourceType) // 简单处理，直接加s
	formalObjectKey := path.Join(formalFolder, filename)

	// 直接上传的方式替代Copy操作，避免Copy操作的参数问题
	// 1. 从临时位置下载文件
	tempFile, err := s.client.Object.Get(context.Background(), tempObjectKey, nil)
	if err != nil {
		return "", fmt.Errorf("获取临时文件失败: %v", err)
	}
	defer tempFile.Body.Close()

	// 2. 将文件重新上传到正式位置
	_, err = s.client.Object.Put(context.Background(), formalObjectKey, tempFile.Body, nil)
	if err != nil {
		return "", fmt.Errorf("重新上传文件失败: %v", err)
	}

	// 3. 删除临时文件
	_, err = s.client.Object.Delete(context.Background(), tempObjectKey)
	if err != nil {
		// 删除临时文件失败不影响正式文件创建，只记录日志
		fmt.Printf("删除临时文件失败: %v\n", err)
	}

	// 返回正式文件URL，确保只有一个斜杠
	baseURL := strings.TrimRight(s.config.BaseURL, "/")
	formalURL := fmt.Sprintf("%s/%s", baseURL, formalObjectKey)
	return formalURL, nil
}

// CleanupTempFiles 清理过期临时文件
// 参数:
//   - olderThan: 清理超过指定时间的临时文件
//
// 返回:
//   - 清理的文件数量
//   - 错误信息
func (s *CosService) CleanupTempFiles(olderThan time.Duration) (int, error) {
	// 构建临时文件夹前缀
	tempPrefix := "tmp/"

	// 列出所有临时文件
	opt := &cos.BucketGetOptions{
		Prefix: tempPrefix,
	}

	var cleanedCount int
	marker := ""

	for {
		resp, _, err := s.client.Bucket.Get(context.Background(), opt)
		if err != nil {
			return cleanedCount, fmt.Errorf("列出临时文件失败: %v", err)
		}

		for _, content := range resp.Contents {
			// 检查文件修改时间
			lastModified, err := time.Parse(time.RFC3339, content.LastModified)
			if err != nil {
				continue
			}

			// 如果文件超过指定时间，删除
			if time.Since(lastModified) > olderThan {
				_, err = s.client.Object.Delete(context.Background(), content.Key)
				if err != nil {
					fmt.Printf("删除临时文件 %s 失败: %v\n", content.Key, err)
					continue
				}
				cleanedCount++
			}
		}

		// 检查是否还有更多文件
		if !resp.IsTruncated {
			break
		}
		marker = resp.NextMarker
		opt.Marker = marker
	}

	return cleanedCount, nil
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
