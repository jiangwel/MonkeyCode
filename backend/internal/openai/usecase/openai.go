package openai

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"

	"github.com/chaitin/MonkeyCode/backend/ent/types"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/domain"
)

type OpenAIUsecase struct {
	repo      domain.OpenAIRepo
	modelRepo domain.ModelRepo
	cfg       *config.Config
	logger    *slog.Logger
}

func NewOpenAIUsecase(
	cfg *config.Config,
	repo domain.OpenAIRepo,
	modelRepo domain.ModelRepo,
	logger *slog.Logger,
) domain.OpenAIUsecase {
	return &OpenAIUsecase{
		repo:      repo,
		modelRepo: modelRepo,
		cfg:       cfg,
		logger:    logger,
	}
}

func (u *OpenAIUsecase) ModelList(ctx context.Context) (*domain.ModelListResp, error) {
	models, err := u.repo.ModelList(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range models {
		u.logger.DebugContext(ctx, "model", slog.Any("model", v))
	}

	resp := &domain.ModelListResp{
		Object: "list",
		Data: cvt.Iter(models, func(_ int, m *db.Model) *domain.ModelData {
			return cvt.From(m, &domain.ModelData{})
		}),
	}

	return resp, nil
}

func (u *OpenAIUsecase) GetConfig(ctx context.Context, req *domain.ConfigReq) (*domain.ConfigResp, error) {
	apiKey, err := u.repo.GetApiKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	llm, err := u.modelRepo.GetWithCache(ctx, consts.ModelTypeLLM)
	if err != nil {
		return nil, err
	}
	coder, err := u.modelRepo.GetWithCache(ctx, consts.ModelTypeCoder)
	if err != nil {
		return nil, err
	}

	if llm.Parameters == nil {
		llm.Parameters = types.DefaultModelParam()
	}

	t, err := template.New("config").Parse(string(config.ConfigTmpl))
	if err != nil {
		return nil, err
	}

	u.logger.With("param", llm.Parameters).DebugContext(ctx, "get config")
	cnt := bytes.NewBuffer(nil)
	data := map[string]any{
		"apiBase":             req.BaseURL,
		"apikey":              apiKey.Key,
		"chatModel":           llm.ModelName,
		"codeModel":           coder.ModelName,
		"r1Enabled":           llm.Parameters.R1Enabled,
		"maxTokens":           llm.Parameters.MaxTokens,
		"contextWindow":       llm.Parameters.ContextWindow,
		"supportsImages":      llm.Parameters.SupprtImages,
		"supportsComputerUse": llm.Parameters.SupportComputerUse,
		"supportsPromptCache": llm.Parameters.SupportPromptCache,
	}
	if err := t.Execute(cnt, data); err != nil {
		return nil, err
	}

	return &domain.ConfigResp{
		Type:    req.Type,
		Content: cnt.String(),
	}, nil
}
