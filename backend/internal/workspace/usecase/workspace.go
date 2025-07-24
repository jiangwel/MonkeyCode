package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/cli"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
)

type WorkspaceFileUsecase struct {
	repo   domain.WorkspaceFileRepo
	config *config.Config
	logger *slog.Logger
}

func NewWorkspaceFileUsecase(
	repo domain.WorkspaceFileRepo,
	config *config.Config,
	logger *slog.Logger,
) domain.WorkspaceFileUsecase {
	return &WorkspaceFileUsecase{
		repo:   repo,
		config: config,
		logger: logger.With("usecase", "workspace_file"),
	}
}

func (u *WorkspaceFileUsecase) Create(ctx context.Context, req *domain.CreateWorkspaceFileReq) (*domain.WorkspaceFile, error) {
	// 验证和计算哈希
	if req.Hash == "" {
		req.Hash = u.calculateHash(req.Content)
	} else if !u.verifyHash(req.Content, req.Hash) {
		return nil, fmt.Errorf("provided hash does not match content")
	}

	// 计算文件大小
	if req.Size == 0 {
		req.Size = int64(len(req.Content))
	}

	// 推断编程语言
	if req.Language == "" {
		req.Language = u.inferLanguage(req.Path)
	}

	file, err := u.repo.Create(ctx, req)
	if err != nil {
		u.logger.Error("failed to create workspace file", "error", err, "path", req.Path)
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	u.logger.Info("workspace file created", "id", file.ID, "path", req.Path)
	return cvt.From(file, &domain.WorkspaceFile{}), nil
}

func (u *WorkspaceFileUsecase) Update(ctx context.Context, req *domain.UpdateWorkspaceFileReq) (*domain.WorkspaceFile, error) {
	file, err := u.repo.Update(ctx, req.ID, func(up *db.WorkspaceFileUpdateOne) error {
		if req.Content != nil {
			up.SetContent(*req.Content)
			// 更新内容时重新计算哈希和大小
			if req.Hash != nil {
				if !u.verifyHash(*req.Content, *req.Hash) {
					return fmt.Errorf("provided hash does not match content")
				}
				up.SetHash(*req.Hash)
			} else {
				up.SetHash(u.calculateHash(*req.Content))
			}
			if req.Size != nil {
				up.SetSize(*req.Size)
			} else {
				up.SetSize(int64(len(*req.Content)))
			}
		} else if req.Hash != nil {
			up.SetHash(*req.Hash)
		}

		if req.Language != nil {
			up.SetLanguage(*req.Language)
		}

		if req.Size != nil && req.Content == nil {
			up.SetSize(*req.Size)
		}

		return nil
	})
	if err != nil {
		u.logger.Error("failed to update workspace file", "error", err, "id", req.ID)
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	u.logger.Info("workspace file updated", "id", req.ID)
	return cvt.From(file, &domain.WorkspaceFile{}), nil
}

func (u *WorkspaceFileUsecase) Delete(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		u.logger.Error("failed to delete workspace file", "error", err, "id", id)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	u.logger.Info("workspace file deleted", "id", id)
	return nil
}

func (u *WorkspaceFileUsecase) GetByID(ctx context.Context, id string) (*domain.WorkspaceFile, error) {
	file, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return cvt.From(file, &domain.WorkspaceFile{}), nil
}

func (u *WorkspaceFileUsecase) GetAndSave(ctx context.Context, req *domain.SaveAstReq) (error) { 
	results, err := cli.RunParseCLI("parse", "", req.Files...) 
	if err != nil {
		return err 
	} 
	for _, res := range results {
		file, err := u.repo.GetByPath(ctx, req.UserID, req.ProjectID, res.FilePath) 
		if err != nil {
			return err 
		} 

		astData, err := json.Marshal(res.Definition) 
		if err != nil {
			return err 
		} 
		_, err = u.repo.Update(ctx, file.ID.String(), func(up *db.WorkspaceFileUpdateOne) error {
			return up.SetContent(string(astData)).Exec(ctx)
		})
		if err != nil {
			return err 
		}
	}
	return nil 
}

func (u *WorkspaceFileUsecase) GetByPath(ctx context.Context, userID, workspaceID, path string) (*domain.WorkspaceFile, error) {
	file, err := u.repo.GetByPath(ctx, userID, workspaceID, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file by path: %w", err)
	}

	return cvt.From(file, &domain.WorkspaceFile{}), nil
}

func (u *WorkspaceFileUsecase) List(ctx context.Context, req *domain.ListWorkspaceFileReq) (*domain.ListWorkspaceFileResp, error) {
	files, pageInfo, err := u.repo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return &domain.ListWorkspaceFileResp{
		PageInfo: pageInfo,
		Files:    domain.FromWorkspaceFiles(files),
	}, nil
}

func (u *WorkspaceFileUsecase) BatchCreate(ctx context.Context, req *domain.BatchCreateWorkspaceFileReq) ([]*domain.WorkspaceFile, error) {
	// 验证和预处理文件
	for _, file := range req.Files {
		if file.Hash == "" {
			file.Hash = u.calculateHash(file.Content)
		} else if !u.verifyHash(file.Content, file.Hash) {
			return nil, fmt.Errorf("hash mismatch for file %s", file.Path)
		}

		if file.Size == 0 {
			file.Size = int64(len(file.Content))
		}

		if file.Language == "" {
			file.Language = u.inferLanguage(file.Path)
		}

		// 设置用户ID和工作区ID
		file.UserID = req.UserID
		file.WorkspaceID = req.WorkspaceID
	}

	files, err := u.repo.BatchCreate(ctx, req.Files)
	if err != nil {
		u.logger.Error("failed to batch create workspace files", "error", err, "count", len(req.Files))
		return nil, fmt.Errorf("failed to batch create files: %w", err)
	}

	u.logger.Info("workspace files batch created", "count", len(files))
	return domain.FromWorkspaceFiles(files), nil
}

func (u *WorkspaceFileUsecase) BatchUpdate(ctx context.Context, req *domain.BatchUpdateWorkspaceFileReq) ([]*domain.WorkspaceFile, error) {
	var results []*domain.WorkspaceFile

	for _, updateReq := range req.Files {
		file, err := u.Update(ctx, updateReq)
		if err != nil {
			return nil, fmt.Errorf("failed to update file %s: %w", updateReq.ID, err)
		}
		results = append(results, file)
	}

	u.logger.Info("workspace files batch updated", "count", len(results))
	return results, nil
}

func (u *WorkspaceFileUsecase) Sync(ctx context.Context, req *domain.SyncWorkspaceFileReq) (*domain.SyncWorkspaceFileResp, error) {
	// 获取要同步的文件哈希列表
	var hashes []string
	fileMap := make(map[string]*domain.CreateWorkspaceFileReq)
	for _, file := range req.Files {
		if file.Hash == "" {
			file.Hash = u.calculateHash(file.Content)
		}
		hashes = append(hashes, file.Hash)
		fileMap[file.Hash] = file
	}

	// 查找工作区中已存在的文件
	existing, err := u.repo.GetByHashes(ctx, req.WorkspaceID, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing files: %w", err)
	}

	var toCreate []*domain.CreateWorkspaceFileReq
	var toUpdate []*domain.UpdateWorkspaceFileReq

	// 分类处理：创建新文件或更新现有文件
	for hash, file := range fileMap {
		file.UserID = req.UserID
		file.WorkspaceID = req.WorkspaceID

		if existingFile, exists := existing[hash]; exists {
			// 文件已存在，检查是否需要更新
			if existingFile.Path != file.Path || existingFile.Language != file.Language {
				updateReq := &domain.UpdateWorkspaceFileReq{
					ID: existingFile.ID.String(),
				}
				if existingFile.Path != file.Path {
					updateReq.Content = &file.Content
				}
				if existingFile.Language != file.Language {
					updateReq.Language = &file.Language
				}
				toUpdate = append(toUpdate, updateReq)
			}
		} else {
			// 新文件，需要创建
			if file.Language == "" {
				file.Language = u.inferLanguage(file.Path)
			}
			if file.Size == 0 {
				file.Size = int64(len(file.Content))
			}
			toCreate = append(toCreate, file)
		}
	}

	resp := &domain.SyncWorkspaceFileResp{}

	// 批量创建新文件
	if len(toCreate) > 0 {
		created, err := u.repo.BatchCreate(ctx, toCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to create new files: %w", err)
		}
		resp.Created = domain.FromWorkspaceFiles(created)
	}

	// 批量更新现有文件
	if len(toUpdate) > 0 {
		updated, err := u.BatchUpdate(ctx, &domain.BatchUpdateWorkspaceFileReq{Files: toUpdate})
		if err != nil {
			return nil, fmt.Errorf("failed to update existing files: %w", err)
		}
		resp.Updated = updated
	}

	resp.Total = len(resp.Created) + len(resp.Updated)

	u.logger.Info("workspace files synced",
		"created", len(resp.Created),
		"updated", len(resp.Updated),
		"total", resp.Total)

	return resp, nil
}

// 辅助方法

func (u *WorkspaceFileUsecase) calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (u *WorkspaceFileUsecase) verifyHash(content, expectedHash string) bool {
	actualHash := u.calculateHash(content)
	return actualHash == expectedHash
}

func (u *WorkspaceFileUsecase) inferLanguage(path string) string {
	// 简单的文件扩展名到语言的映射
	if idx := strings.LastIndex(path, "."); idx != -1 {
		ext := strings.ToLower(path[idx+1:])
		switch ext {
		case "go":
			return "go"
		case "js", "mjs":
			return "javascript"
		case "ts":
			return "typescript"
		case "py":
			return "python"
		case "java":
			return "java"
		case "cpp", "cc", "cxx":
			return "cpp"
		case "c":
			return "c"
		case "rs":
			return "rust"
		case "php":
			return "php"
		case "rb":
			return "ruby"
		case "swift":
			return "swift"
		case "kt":
			return "kotlin"
		case "cs":
			return "csharp"
		case "sh", "bash":
			return "shell"
		case "sql":
			return "sql"
		case "yaml", "yml":
			return "yaml"
		case "json":
			return "json"
		case "xml":
			return "xml"
		case "html":
			return "html"
		case "css":
			return "css"
		case "md":
			return "markdown"
		case "toml":
			return "toml"
		case "ini":
			return "ini"
		default:
			return "text"
		}
	}
	return "text"
}
