package domain

import (
	"context"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
)

// WorkspaceFileUsecase 定义 WorkspaceFile 业务逻辑接口
type WorkspaceFileUsecase interface {
	Create(ctx context.Context, req *CreateWorkspaceFileReq) (*WorkspaceFile, error)
	Update(ctx context.Context, req *UpdateWorkspaceFileReq) (*WorkspaceFile, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*WorkspaceFile, error)
	GetByPath(ctx context.Context, userID, workspaceID, path string) (*WorkspaceFile, error)
	List(ctx context.Context, req *ListWorkspaceFileReq) (*ListWorkspaceFileResp, error)
	BatchCreate(ctx context.Context, req *BatchCreateWorkspaceFileReq) ([]*WorkspaceFile, error)
	BatchUpdate(ctx context.Context, req *BatchUpdateWorkspaceFileReq) ([]*WorkspaceFile, error)
	Sync(ctx context.Context, req *SyncWorkspaceFileReq) (*SyncWorkspaceFileResp, error)
}

// WorkspaceFileRepo 定义 WorkspaceFile 数据访问接口
type WorkspaceFileRepo interface {
	Create(ctx context.Context, req *CreateWorkspaceFileReq) (*db.WorkspaceFile, error)
	Update(ctx context.Context, id string, fn func(*db.WorkspaceFileUpdateOne) error) (*db.WorkspaceFile, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*db.WorkspaceFile, error)
	GetByPath(ctx context.Context, userID, workspaceID, path string) (*db.WorkspaceFile, error)
	List(ctx context.Context, req *ListWorkspaceFileReq) ([]*db.WorkspaceFile, *db.PageInfo, error)
	BatchCreate(ctx context.Context, files []*CreateWorkspaceFileReq) ([]*db.WorkspaceFile, error)
	GetByHashes(ctx context.Context, workspaceID string, hashes []string) (map[string]*db.WorkspaceFile, error)
	CountByWorkspace(ctx context.Context, workspaceID string) (int64, error)
	GetWorkspaceFiles(ctx context.Context, workspaceID string) ([]*db.WorkspaceFile, error)
}

// 请求结构体

type CreateWorkspaceFileReq struct {
	UserID      string `json:"user_id" validate:"required"`      // 用户ID
	WorkspaceID string `json:"workspace_id" validate:"required"` // 工作区ID
	Path        string `json:"path" validate:"required"`         // 文件路径
	Content     string `json:"content"`                          // 文件内容
	Hash        string `json:"hash" validate:"required"`         // 文件哈希
	Language    string `json:"language"`                         // 编程语言
	Size        int64  `json:"size"`                             // 文件大小
}

type UpdateWorkspaceFileReq struct {
	ID       string  `json:"id" validate:"required"` // 文件ID
	Content  *string `json:"content"`                // 文件内容
	Hash     *string `json:"hash"`                   // 文件哈希
	Language *string `json:"language"`               // 编程语言
	Size     *int64  `json:"size"`                   // 文件大小
}

type ListWorkspaceFileReq struct {
	*web.Pagination
	UserID      string `json:"user_id" query:"user_id"`           // 用户ID
	WorkspaceID string `json:"workspace_id" query:"workspace_id"` // 工作区ID
	Language    string `json:"language" query:"language"`         // 编程语言筛选
	Search      string `json:"search" query:"search"`             // 搜索关键词（文件路径）
}

type BatchCreateWorkspaceFileReq struct {
	UserID      string                    `json:"user_id" validate:"required"`      // 用户ID
	WorkspaceID string                    `json:"workspace_id" validate:"required"` // 工作区ID
	Files       []*CreateWorkspaceFileReq `json:"files" validate:"required,dive"`   // 文件列表
}

type BatchUpdateWorkspaceFileReq struct {
	Files []*UpdateWorkspaceFileReq `json:"files" validate:"required,dive"` // 文件列表
}

type SyncWorkspaceFileReq struct {
	UserID      string                    `json:"user_id" validate:"required"`      // 用户ID
	WorkspaceID string                    `json:"workspace_id" validate:"required"` // 工作区ID
	Files       []*CreateWorkspaceFileReq `json:"files" validate:"required,dive"`   // 要同步的文件列表
}

// 响应结构体

type ListWorkspaceFileResp struct {
	*db.PageInfo
	Files []*WorkspaceFile `json:"files"`
}

type SyncWorkspaceFileResp struct {
	Created []*WorkspaceFile `json:"created"` // 新创建的文件
	Updated []*WorkspaceFile `json:"updated"` // 更新的文件
	Deleted []string         `json:"deleted"` // 删除的文件ID
	Total   int              `json:"total"`   // 处理的文件总数
}

// 数据模型

type WorkspaceFile struct {
	ID          string `json:"id"`           // 文件ID
	UserID      string `json:"user_id"`      // 用户ID
	WorkspaceID string `json:"workspace_id"` // 工作区ID
	Path        string `json:"path"`         // 文件路径
	Content     string `json:"content"`      // 文件内容
	Hash        string `json:"hash"`         // 文件哈希
	Language    string `json:"language"`     // 编程语言
	Size        int64  `json:"size"`         // 文件大小
	CreatedAt   int64  `json:"created_at"`   // 创建时间
	UpdatedAt   int64  `json:"updated_at"`   // 更新时间
}

func (w *WorkspaceFile) From(e *db.WorkspaceFile) *WorkspaceFile {
	if e == nil {
		return w
	}

	w.ID = e.ID.String()
	w.UserID = e.UserID.String()
	w.WorkspaceID = e.WorkspaceID.String()
	w.Path = e.Path
	w.Content = e.Content
	w.Hash = e.Hash
	w.Language = e.Language
	w.Size = e.Size
	w.CreatedAt = e.CreatedAt.Unix()
	w.UpdatedAt = e.UpdatedAt.Unix()

	return w
}

// 工具函数
func FromWorkspaceFiles(files []*db.WorkspaceFile) []*WorkspaceFile {
	return cvt.Iter(files, func(_ int, e *db.WorkspaceFile) *WorkspaceFile {
		return cvt.From(e, &WorkspaceFile{})
	})
}
