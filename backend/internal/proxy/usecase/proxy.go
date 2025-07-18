package usecase

import (
	"context"

	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"

	"github.com/chaitin/MonkeyCode/backend/consts"
	"github.com/chaitin/MonkeyCode/backend/domain"
)

type ProxyUsecase struct {
	repo      domain.ProxyRepo
	modelRepo domain.ModelRepo
}

func NewProxyUsecase(repo domain.ProxyRepo, modelRepo domain.ModelRepo) domain.ProxyUsecase {
	return &ProxyUsecase{repo: repo, modelRepo: modelRepo}
}

func (p *ProxyUsecase) Record(ctx context.Context, record *domain.RecordParam) error {
	return p.repo.Record(ctx, record)
}

// SelectModelWithLoadBalancing implements domain.ProxyUsecase.
func (p *ProxyUsecase) SelectModelWithLoadBalancing(modelName string, modelType consts.ModelType) (*domain.Model, error) {
	model, err := p.modelRepo.GetWithCache(context.Background(), modelType)
	if err != nil {
		return nil, err
	}
	return cvt.From(model, &domain.Model{}), nil
}

func (p *ProxyUsecase) ValidateApiKey(ctx context.Context, key string) (*domain.ApiKey, error) {
	apiKey, err := p.repo.ValidateApiKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return cvt.From(apiKey, &domain.ApiKey{}), nil
}

func (p *ProxyUsecase) AcceptCompletion(ctx context.Context, req *domain.AcceptCompletionReq) error {
	return p.repo.AcceptCompletion(ctx, req)
}

func (p *ProxyUsecase) Report(ctx context.Context, req *domain.ReportReq) error {
	return p.repo.Report(ctx, req)
}
