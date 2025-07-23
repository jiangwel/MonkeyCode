package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/report"
	"github.com/chaitin/MonkeyCode/backend/pkg/version"
)

type ReportUsecase struct {
	repo     domain.ReportRepo
	logger   *slog.Logger
	reporter *report.Reporter
}

func NewReportUsecase(repo domain.ReportRepo, logger *slog.Logger, reporter *report.Reporter) domain.ReportUsecase {
	r := &ReportUsecase{repo: repo, logger: logger, reporter: reporter}
	go r.Report()
	return r
}

func (r *ReportUsecase) Report() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		if err := r.innerReport(); err != nil {
			r.logger.With("error", err).Error("report failed")
		}
	}
}

func (r *ReportUsecase) innerReport() error {
	admins, err := r.repo.GetAdminCount(context.Background())
	if err != nil {
		return fmt.Errorf("get admin count failed: %w", err)
	}
	models, err := r.repo.GetCurrentModels(context.Background())
	if err != nil {
		return fmt.Errorf("get current models failed: %w", err)
	}
	stats, err := r.repo.GetLast24HoursStats(context.Background())
	if err != nil {
		return fmt.Errorf("get last 24 hours stats failed: %w", err)
	}
	members, err := r.repo.GetMemberCount(context.Background())
	if err != nil {
		return fmt.Errorf("get member count failed: %w", err)
	}

	data := domain.ReportData{
		Timestamp:     time.Now().Format(time.RFC3339),
		Version:       version.Version,
		MachineID:     r.reporter.GetMachineID(),
		AdminCount:    admins,
		MemberCount:   members,
		Last24Hours:   stats,
		CurrentModels: models,
	}

	if err := r.reporter.Report("monkeycode-metrics", data); err != nil {
		return fmt.Errorf("report failed: %w", err)
	}
	return nil
}
