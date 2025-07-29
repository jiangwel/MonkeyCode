package v1

import (
	"fmt"
	"log/slog"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/internal/middleware"
	"github.com/chaitin/MonkeyCode/backend/pkg/logger"
)

type CodeSnippetHandler struct {
	usecase domain.CodeSnippetUsecase
	logger  *slog.Logger
}

func NewCodeSnippetHandler(
	w *web.Web,
	usecase domain.CodeSnippetUsecase,
	auth *middleware.AuthMiddleware,
	active *middleware.ActiveMiddleware,
	readonly *middleware.ReadOnlyMiddleware,
	proxy *middleware.ProxyMiddleware,
	logger *slog.Logger,
) *CodeSnippetHandler {
	h := &CodeSnippetHandler{
		usecase: usecase,
		logger:  logger.With("handler", "codesnippet"),
	}

	// 设置路由 - 使用API Key认证的接口（IDE端使用）
	ide := w.Group("/api/v1/ide/codesnippet")
	ide.Use(proxy.Auth(), active.Active("apikey"), readonly.Guard())

	// IDE端上下文检索接口
	ide.POST("/context", web.BindHandler(h.GetContext))

	return h
}

// GetContextReq IDE端上下文检索请求
type GetContextReq struct {
	// 批量查询参数
	Queries []Query `json:"queries,omitempty"` // 批量查询条件

	// 单个查询参数
	Query Query `json:"query,omitempty"` // 单个查询条件

	Limit         int    `json:"limit"`          // 返回结果数量限制，默认10
	WorkspacePath string `json:"workspace_path"` // 工作区路径（必填）
}

// Query 批量查询条件
type Query struct {
	Name     string `json:"name,omitempty"`     // 代码片段名称（可选）
	Type     string `json:"type,omitempty"`     // 代码片段类型（可选）
	Language string `json:"language,omitempty"` // 编程语言（可选）
}

// GetContext IDE端上下文检索接口
//
//	@Tags			CodeSnippet
//	@Summary		IDE端上下文检索
//	@Description	为IDE端提供代码片段上下文检索功能，使用API Key认证。支持单个查询和批量查询。
//	@ID				get-context
//	@Accept			json
//	@Produce		json
//	@Param			request	body		GetContextReq	true	"检索请求参数"
//	@Success		200		{object}	web.Resp{data=[]domain.CodeSnippet}
//	@Router			/api/v1/ide/codesnippet/context [post]
//	@Security		ApiKeyAuth
func (h *CodeSnippetHandler) GetContext(c *web.Context, req GetContextReq) error {
	// 设置默认限制
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50 // 最大限制50个结果
	}

	// 如果没有提供workspace_path，则返回空结果
	if req.WorkspacePath == "" {
		return c.Success([]*domain.CodeSnippet{})
	}

	// 获取用户ID，主要使用API Key认证
	var userID string
	if ctxUserID := c.Request().Context().Value(logger.UserIDKey{}); ctxUserID != nil {
		userID = ctxUserID.(string)
	} else {
		h.logger.Error("API Key authentication required for IDE context retrieval")
		return fmt.Errorf("API Key authentication required")
	}

	var allSnippets []*domain.CodeSnippet

	// 如果提供了批量查询条件，则执行批量查询
	if len(req.Queries) > 0 {
		// 执行批量查询
		for _, query := range req.Queries {
			snippets, err := h.usecase.SearchByWorkspace(c.Request().Context(), userID, req.WorkspacePath, query.Name, query.Type, query.Language)
			if err != nil {
				h.logger.Error("failed to get context for IDE", "error", err)
				return err
			}
			allSnippets = append(allSnippets, snippets...)
		}
	} else {
		// 执行单个查询
		snippets, err := h.usecase.SearchByWorkspace(c.Request().Context(), userID, req.WorkspacePath, req.Query.Name, req.Query.Type, req.Query.Language)
		if err != nil {
			h.logger.Error("failed to get context for IDE", "error", err)
			return err
		}
		allSnippets = snippets
	}

	// 限制返回结果数量
	if len(allSnippets) > req.Limit {
		allSnippets = allSnippets[:req.Limit]
	}

	h.logger.Info("IDE context retrieval completed",
		"userID", c.Request().Context().Value(logger.UserIDKey{}),
		"resultCount", len(allSnippets),
		"filters", map[string]interface{}{
			"singleQuery": map[string]string{
				"name":     req.Query.Name,
				"type":     req.Query.Type,
				"language": req.Query.Language,
			},
			"batchQueryCount": len(req.Queries),
			"workspace_path":  req.WorkspacePath,
		})

	return c.Success(allSnippets)
}
