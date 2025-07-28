package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/google/uuid"
)

type CodeSnippetRepo struct {
	client *db.Client
	logger *slog.Logger
}

func NewCodeSnippetRepo(client *db.Client, logger *slog.Logger) domain.CodeSnippetRepo {
	return &CodeSnippetRepo{
		client: client,
		logger: logger.With("repo", "codesnippet"),
	}
}

func (r *CodeSnippetRepo) Create(ctx context.Context, req *domain.CreateCodeSnippetReq) (*db.CodeSnippet, error) {
	// 将 workspaceFileID 字符串转换为 UUID
	workspaceFileUUID, err := uuid.Parse(req.WorkspaceFileID)
	if err != nil {
		r.logger.Error("failed to parse workspace file ID", "error", err, "id", req.WorkspaceFileID)
		return nil, fmt.Errorf("invalid workspace file ID: %w", err)
	}

	create := r.client.CodeSnippet.Create().
		SetWorkspaceFileID(workspaceFileUUID).
		SetName(req.Name).
		SetSnippetType(req.SnippetType).
		SetLanguage(req.Language).
		SetContent(req.Content).
		SetHash(req.Hash).
		SetStartLine(req.StartLine).
		SetEndLine(req.EndLine).
		SetStartColumn(req.StartColumn).
		SetEndColumn(req.EndColumn).
		SetNamespace(req.Namespace).
		SetContainerName(req.ContainerName).
		SetScope(req.Scope).
		SetDependencies(req.Dependencies).
		SetParameters(req.Parameters).
		SetSignature(req.Signature).
		SetDefinitionText(req.DefinitionText).
		SetStructuredInfo(req.StructuredInfo)

	snippet, err := create.Save(ctx)
	if err != nil {
		r.logger.Error("failed to create code snippet", "error", err)
		return nil, err
	}

	return snippet, nil
}

func (r *CodeSnippetRepo) ListByWorkspaceFile(ctx context.Context, workspaceFileID string) ([]*db.CodeSnippet, error) {
	// 实现列出特定工作区文件的所有代码片段的逻辑
	// 这里需要将 workspaceFileID 字符串转换为 UUID
	// 为简化起见，这里暂时返回空列表，实际实现需要根据需求完成
	return []*db.CodeSnippet{}, nil
}

func (r *CodeSnippetRepo) GetByID(ctx context.Context, id string) (*db.CodeSnippet, error) {
	// 实现根据 ID 获取代码片段的逻辑
	// 这里需要将 id 字符串转换为 UUID
	// 为简化起见，这里暂时返回 nil，实际实现需要根据需求完成
	return nil, nil
}

func (r *CodeSnippetRepo) Delete(ctx context.Context, id string) error {
	// 实现删除代码片段的逻辑
	// 这里需要将 id 字符串转换为 UUID
	// 为简化起见，这里暂时返回 nil，实际实现需要根据需求完成
	return nil
}
