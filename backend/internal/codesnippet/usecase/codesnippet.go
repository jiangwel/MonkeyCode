package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chaitin/MonkeyCode/backend/domain"
)

type CodeSnippetUsecase struct {
	repo   domain.CodeSnippetRepo
	logger *slog.Logger
}

func NewCodeSnippetUsecase(
	repo domain.CodeSnippetRepo,
	logger *slog.Logger,
) domain.CodeSnippetUsecase {
	return &CodeSnippetUsecase{
		repo:   repo,
		logger: logger.With("usecase", "codesnippet"),
	}
}

// CreateFromIndexResult 从 IndexResult 创建 CodeSnippet
func (u *CodeSnippetUsecase) CreateFromIndexResult(ctx context.Context, workspaceFileID string, indexResult *domain.IndexResult) (*domain.CodeSnippet, error) {
	// 构建 CreateCodeSnippetReq
	req := &domain.CreateCodeSnippetReq{
		WorkspaceFileID: workspaceFileID,
		Name:            indexResult.Name,
		SnippetType:     indexResult.Type,
		Language:        indexResult.Language,
		Content:         indexResult.RangeText,
		Hash:            indexResult.FileHash,
		StartLine:       indexResult.StartLine,
		EndLine:         indexResult.EndLine,
		// StartColumn 和 EndColumn 在 IndexResult 中没有直接对应字段，暂时设置为 0
		StartColumn:    0,
		EndColumn:      0,
		Namespace:      "",                 // IndexResult 中没有直接对应字段
		ContainerName:  "",                 // IndexResult 中没有直接对应字段
		Dependencies:   []string{},         // IndexResult 中没有直接对应字段
		Parameters:     []map[string]any{}, // IndexResult 中没有直接对应字段
		Signature:      indexResult.Signature,
		DefinitionText: indexResult.DefinitionText,
		StructuredInfo: map[string]any{
			"definition": indexResult.Definition,
		},
	}

	// 创建 CodeSnippet
	snippet, err := u.repo.Create(ctx, req)
	if err != nil {
		u.logger.Error("failed to create code snippet from index result", "error", err)
		return nil, fmt.Errorf("failed to create code snippet: %w", err)
	}

	// 转换为领域模型
	return (&domain.CodeSnippet{}).From(snippet), nil
}

// ListByWorkspaceFile 列出特定工作区文件的所有代码片段
func (u *CodeSnippetUsecase) ListByWorkspaceFile(ctx context.Context, workspaceFileID string) ([]*domain.CodeSnippet, error) {
	// 实现列出特定工作区文件的所有代码片段的逻辑
	// 为简化起见，这里暂时返回空列表，实际实现需要根据需求完成
	return []*domain.CodeSnippet{}, nil
}

// GetByID 根据 ID 获取代码片段
func (u *CodeSnippetUsecase) GetByID(ctx context.Context, id string) (*domain.CodeSnippet, error) {
	// 实现根据 ID 获取代码片段的逻辑
	// 为简化起见，这里暂时返回 nil，实际实现需要根据需求完成
	return nil, nil
}

// Delete 删除代码片段
func (u *CodeSnippetUsecase) Delete(ctx context.Context, id string) error {
	// 实现删除代码片段的逻辑
	// 为简化起见，这里暂时返回 nil，实际实现需要根据需求完成
	return nil
}
