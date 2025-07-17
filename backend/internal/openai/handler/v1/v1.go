package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/internal/middleware"
	"github.com/chaitin/MonkeyCode/backend/internal/proxy"
)

type V1Handler struct {
	logger   *slog.Logger
	proxy    *proxy.LLMProxy
	proxyUse domain.ProxyUsecase
	usecase  domain.OpenAIUsecase
	euse     domain.ExtensionUsecase
	config   *config.Config
}

func NewV1Handler(
	logger *slog.Logger,
	w *web.Web,
	proxy *proxy.LLMProxy,
	proxyUse domain.ProxyUsecase,
	usecase domain.OpenAIUsecase,
	euse domain.ExtensionUsecase,
	middleware *middleware.ProxyMiddleware,
	active *middleware.ActiveMiddleware,
	config *config.Config,
) *V1Handler {
	h := &V1Handler{
		logger:   logger.With(slog.String("handler", "openai")),
		proxy:    proxy,
		proxyUse: proxyUse,
		usecase:  usecase,
		euse:     euse,
		config:   config,
	}

	w.GET("/api/config", web.BindHandler(h.GetConfig), middleware.Auth())
	w.GET("/v1/version", web.BaseHandler(h.Version), middleware.Auth())

	g := w.Group("/v1", middleware.Auth())
	g.GET("/models", web.BaseHandler(h.ModelList))
	g.POST("/completion/accept", web.BindHandler(h.AcceptCompletion), active.Active("user"))
	g.POST("/chat/completions", web.BaseHandler(h.ChatCompletion), active.Active("user"))
	g.POST("/completions", web.BaseHandler(h.Completions), active.Active("user"))
	g.POST("/embeddings", web.BaseHandler(h.Embeddings), active.Active("user"))
	return h
}

func BadRequest(c *web.Context, msg string) error {
	c.JSON(http.StatusBadRequest, echo.Map{
		"error": echo.Map{
			"message": msg,
			"type":    "invalid_request_error",
		},
	})
	return nil
}

func (h *V1Handler) Version(c *web.Context) error {
	v, err := h.euse.Latest(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, domain.VersionInfo{
		Version: v.Version,
		URL:     fmt.Sprintf("%s/api/v1/static/vsix/%s", h.config.BaseUrl, v.Version),
	})
}

// AcceptCompletion 接受补全请求
//
//	@Tags			OpenAIV1
//	@Summary		接受补全请求
//	@Description	接受补全请求
//	@ID				accept-completion
//	@Accept			json
//	@Produce		json
//	@Param			param	body		domain.AcceptCompletionReq	true	"补全请求"
//	@Success		200		{object}	web.Resp{}
//	@Router			/v1/completion/accept [post]
func (h *V1Handler) AcceptCompletion(c *web.Context, req domain.AcceptCompletionReq) error {
	if err := h.proxyUse.AcceptCompletion(c.Request().Context(), &req); err != nil {
		return BadRequest(c, err.Error())
	}
	return nil
}

// ModelList 模型列表
//
//	@Tags			OpenAIV1
//	@Summary		模型列表
//	@Description	模型列表
//	@ID				model-list
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{data=domain.ModelListResp}
//	@Router			/v1/models [get]
func (h *V1Handler) ModelList(c *web.Context) error {
	resp, err := h.usecase.ModelList(c.Request().Context())
	if err != nil {
		return err
	}
	return c.Success(resp)
}

// ChatCompletion 处理聊天补全请求
//
//	@Tags			OpenAIV1
//	@Summary		处理聊天补全请求
//	@Description	处理聊天补全请求
//	@ID				chat-completion
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/chat/completions [post]
func (h *V1Handler) ChatCompletion(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Completions 处理文本补全请求
//
//	@Tags			OpenAIV1
//	@Summary		处理文本补全请求
//	@Description	处理文本补全请求
//	@ID				completions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/completions [post]
func (h *V1Handler) Completions(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Embeddings 处理嵌入请求
//
//	@Tags			OpenAIV1
//	@Summary		处理嵌入请求
//	@Description	处理嵌入请求
//	@ID				embeddings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/embeddings [post]
func (h *V1Handler) Embeddings(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *V1Handler) GetConfig(c *web.Context, req domain.ConfigReq) error {
	key := middleware.GetApiKey(c)
	req.Key = key.Key
	resp, err := h.usecase.GetConfig(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
