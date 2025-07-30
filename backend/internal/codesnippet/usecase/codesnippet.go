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
	// 调用 repository 层的方法
	dbSnippets, err := u.repo.ListByWorkspaceFile(ctx, workspaceFileID)
	if err != nil {
		u.logger.Error("failed to list code snippets by workspace file", "error", err, "workspaceFileID", workspaceFileID)
		return nil, fmt.Errorf("failed to list code snippets: %w", err)
	}

	// 将数据库模型转换为领域模型
	var snippets []*domain.CodeSnippet
	for _, dbSnippet := range dbSnippets {
		snippet := (&domain.CodeSnippet{}).From(dbSnippet)
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

// GetByID 根据 ID 获取代码片段
func (u *CodeSnippetUsecase) GetByID(ctx context.Context, id string) (*domain.CodeSnippet, error) {
	// 调用 repository 层的方法
	dbSnippet, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("failed to get code snippet by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get code snippet: %w", err)
	}

	// 将数据库模型转换为领域模型
	return (&domain.CodeSnippet{}).From(dbSnippet), nil
}

// Delete 删除代码片段
func (u *CodeSnippetUsecase) Delete(ctx context.Context, id string) error {
	// 调用 repository 层的方法
	err := u.repo.Delete(ctx, id)
	if err != nil {
		u.logger.Error("failed to delete code snippet", "error", err, "id", id)
		return fmt.Errorf("failed to delete code snippet: %w", err)
	}

	return nil
}

// Search 根据名称、类型和语言搜索代码片段
func (u *CodeSnippetUsecase) Search(ctx context.Context, name, snippetType, language string) ([]*domain.CodeSnippet, error) {
	// 调用 repository 层的 Search 方法
	dbSnippets, err := u.repo.Search(ctx, name, snippetType, language)
	if err != nil {
		u.logger.Error("failed to search code snippets", "error", err)
		return nil, fmt.Errorf("failed to search code snippets: %w", err)
	}

	// 将数据库模型转换为领域模型
	var snippets []*domain.CodeSnippet
	for _, dbSnippet := range dbSnippets {
		snippet := (&domain.CodeSnippet{}).From(dbSnippet)
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

// SearchByWorkspace 根据用户ID、工作区路径和搜索条件搜索代码片段
func (u *CodeSnippetUsecase) SearchByWorkspace(ctx context.Context, userID, workspacePath, name, snippetType, language string) ([]*domain.CodeSnippet, error) {
	// 调用 repository 层的 SearchByWorkspace 方法
	dbSnippets, err := u.repo.SearchByWorkspace(ctx, userID, workspacePath, name, snippetType, language)
	if err != nil {
		u.logger.Error("failed to search code snippets by workspace", "error", err, "userID", userID, "workspacePath", workspacePath)
		return nil, fmt.Errorf("failed to search code snippets by workspace: %w", err)
	}

	// 将数据库模型转换为领域模型
	var snippets []*domain.CodeSnippet
	for _, dbSnippet := range dbSnippets {
		snippet := (&domain.CodeSnippet{}).FromWithFile(dbSnippet)
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}
