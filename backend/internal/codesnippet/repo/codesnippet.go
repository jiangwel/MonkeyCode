package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/db/codesnippet"
	"github.com/chaitin/MonkeyCode/backend/db/predicate"
	"github.com/chaitin/MonkeyCode/backend/db/workspace"
	"github.com/chaitin/MonkeyCode/backend/db/workspacefile"
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
	// 将 workspaceFileID 字符串转换为 UUID
	workspaceFileUUID, err := uuid.Parse(workspaceFileID)
	if err != nil {
		r.logger.Error("failed to parse workspace file ID", "error", err, "id", workspaceFileID)
		return nil, fmt.Errorf("invalid workspace file ID: %w", err)
	}

	// 查询特定工作区文件的所有代码片段
	snippets, err := r.client.CodeSnippet.Query().
		Where(codesnippet.WorkspaceFileID(workspaceFileUUID)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to list code snippets by workspace file", "error", err, "workspaceFileID", workspaceFileID)
		return nil, fmt.Errorf("failed to list code snippets: %w", err)
	}

	return snippets, nil
}

func (r *CodeSnippetRepo) GetByID(ctx context.Context, id string) (*db.CodeSnippet, error) {
	// 检查ID是否为空
	if id == "" {
		r.logger.Error("code snippet ID is empty", "id", id)
		return nil, fmt.Errorf("invalid code snippet ID: empty ID")
	}

	// 将 id 字符串转换为 UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("failed to parse code snippet ID", "error", err, "id", id)
		return nil, fmt.Errorf("invalid code snippet ID: %w", err)
	}

	// 根据 ID 获取代码片段
	snippet, err := r.client.CodeSnippet.Get(ctx, uuid)
	if err != nil {
		r.logger.Error("failed to get code snippet by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get code snippet: %w", err)
	}

	return snippet, nil
}

func (r *CodeSnippetRepo) Delete(ctx context.Context, id string) error {
	// 检查ID是否为空
	if id == "" {
		r.logger.Error("code snippet ID is empty", "id", id)
		return fmt.Errorf("invalid code snippet ID: empty ID")
	}

	// 将 id 字符串转换为 UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		r.logger.Error("failed to parse code snippet ID", "error", err, "id", id)
		return fmt.Errorf("invalid code snippet ID: %w", err)
	}

	// 删除代码片段
	err = r.client.CodeSnippet.DeleteOneID(uuid).Exec(ctx)
	if err != nil {
		r.logger.Error("failed to delete code snippet", "error", err, "id", id)
		return fmt.Errorf("failed to delete code snippet: %w", err)
	}

	return nil
}

// Search 根据名称、类型和语言搜索代码片段
func (r *CodeSnippetRepo) Search(ctx context.Context, name, snippetType, language string) ([]*db.CodeSnippet, error) {
	// 构建查询
	query := r.client.CodeSnippet.Query()

	// 创建一个切片来存储所有谓词
	var predicates []predicate.CodeSnippet

	// 如果提供了名称参数，则添加名称过滤条件
	if name != "" {
		predicates = append(predicates, codesnippet.Name(name))
	}

	// 如果提供了类型参数，则添加类型过滤条件
	if snippetType != "" {
		predicates = append(predicates, codesnippet.SnippetType(snippetType))
	}

	// 如果提供了语言参数，则添加语言过滤条件
	if language != "" {
		predicates = append(predicates, codesnippet.Language(language))
	}

	// 如果有任何谓词，将它们添加到查询中
	if len(predicates) > 0 {
		query = query.Where(codesnippet.And(predicates...))
	}

	// 执行查询并返回结果
	return query.All(ctx)
}

// SearchByWorkspace 根据用户ID、工作区路径和搜索条件搜索代码片段
// 只有在提供了至少一个搜索条件时才返回结果，否则返回空数组
func (r *CodeSnippetRepo) SearchByWorkspace(ctx context.Context, userID, workspacePath, name, snippetType, language string) ([]*db.CodeSnippet, error) {
	// 检查是否提供了至少一个搜索条件
	if name == "" && snippetType == "" && language == "" {
		// 如果没有提供任何搜索条件，返回空结果
		return []*db.CodeSnippet{}, nil
	}

	// 首先根据用户ID和工作区路径找到工作区
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		r.logger.Error("failed to parse user ID", "error", err, "userID", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// 查询工作区
	workspace, err := r.client.Workspace.Query().
		Where(
			workspace.UserID(userUUID),
			workspace.RootPath(workspacePath),
		).
		Only(ctx)
	if err != nil {
		r.logger.Error("failed to find workspace", "error", err, "userID", userID, "workspacePath", workspacePath)
		return nil, fmt.Errorf("workspace not found: %w", err)
	}

	// 查询该工作区下的所有文件
	workspaceFiles, err := r.client.WorkspaceFile.Query().
		Where(workspacefile.WorkspaceID(workspace.ID)).
		All(ctx)
	if err != nil {
		r.logger.Error("failed to get workspace files", "error", err, "workspaceID", workspace.ID)
		return nil, fmt.Errorf("failed to get workspace files: %w", err)
	}

	if len(workspaceFiles) == 0 {
		return []*db.CodeSnippet{}, nil
	}

	// 提取文件ID列表
	var fileIDs []uuid.UUID
	for _, file := range workspaceFiles {
		fileIDs = append(fileIDs, file.ID)
	}

	// 构建代码片段查询，包含WorkspaceFile信息
	query := r.client.CodeSnippet.Query().
		Where(codesnippet.WorkspaceFileIDIn(fileIDs...)).
		WithSourceFile() // 预加载WorkspaceFile信息

	// 创建一个切片来存储所有谓词
	var predicates []predicate.CodeSnippet

	// 如果提供了名称参数，则添加名称过滤条件
	if name != "" {
		predicates = append(predicates, codesnippet.Name(name))
	}

	// 如果提供了类型参数，则添加类型过滤条件
	if snippetType != "" {
		predicates = append(predicates, codesnippet.SnippetType(snippetType))
	}

	// 如果提供了语言参数，则添加语言过滤条件
	if language != "" {
		predicates = append(predicates, codesnippet.Language(language))
	}

	// 添加谓词到查询中（这里总是会添加，因为我们已经检查了至少有一个条件）
	query = query.Where(codesnippet.And(predicates...))

	// 执行查询并返回结果
	return query.All(ctx)
}
