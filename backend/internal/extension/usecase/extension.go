package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/db"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/pkg/cvt"
)

type ExtensionUsecase struct {
	logger  *slog.Logger
	repo    domain.ExtensionRepo
	config  *config.Config
	version string
	mu      sync.Mutex
}

func NewExtensionUsecase(
	repo domain.ExtensionRepo,
	config *config.Config,
	logger *slog.Logger,
) domain.ExtensionUsecase {
	e := &ExtensionUsecase{
		repo:   repo,
		config: config,
		mu:     sync.Mutex{},
		logger: logger,
	}
	go e.syncLatest()
	return e
}

// GetByVersion implements domain.ExtensionUsecase.
func (e *ExtensionUsecase) GetByVersion(ctx context.Context, version string) (*domain.Extension, error) {
	ee, err := e.repo.GetByVersion(ctx, version)
	if err != nil {
		return nil, err
	}

	return cvt.From(ee, &domain.Extension{}), nil
}

// Latest implements domain.ExtensionUsecase.
func (e *ExtensionUsecase) Latest(ctx context.Context) (*domain.Extension, error) {
	ee, err := e.repo.Latest(ctx)
	if err != nil {
		return nil, err
	}

	return cvt.From(ee, &domain.Extension{}), nil
}

func (e *ExtensionUsecase) syncLatest() {
	latest, err := e.repo.Latest(context.Background())
	if err == nil {
		e.version = latest.Version
	}

	logger := e.logger.With("fn", "syncLatest")
	logger.With("version", e.version).Debug("开始同步插件版本信息")
	t := time.NewTicker(1 * time.Minute)

	e.innerSync()
	for range t.C {
		e.innerSync()
	}
}

func (e *ExtensionUsecase) innerSync() {
	logger := e.logger.With("fn", "syncLatest")
	v, err := e.getRemoteVersion()
	if err != nil {
		logger.With("error", err).Error("获取远程插件版本信息失败")
		return
	}
	if v == e.version {
		return
	}

	e.download(v)
}

func (e *ExtensionUsecase) download(version string) {
	logger := e.logger.With("fn", "download")
	url := fmt.Sprintf("%s/monkeycode/vsixs/monkeycode-%s.vsix", e.config.Extension.Baseurl, version)
	logger.With("url", url).With("version", version).Debug("发现新版本，开始下载")
	resp, err := http.Get(url)
	if err != nil {
		logger.With("error", err).Error("下载插件失败")
		return
	}
	defer resp.Body.Close()
	filename := fmt.Sprintf("/app/static/monkeycode-%s.vsix", version)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.With("error", err).Error("创建文件失败")
		return
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		logger.With("error", err).Error("复制文件内容失败")
		return
	}
	logger.Debug("下载插件成功")

	if _, err := e.repo.Save(context.Background(), &db.Extension{
		Version: version,
		Path:    filename,
	}); err != nil {
		logger.With("error", err).Error("保存插件版本信息失败")
		os.Remove(filename)
		return
	}
	e.version = version
}

func (e *ExtensionUsecase) getRemoteVersion() (string, error) {
	logger := e.logger.With("fn", "getRemoteVersion")
	u := fmt.Sprintf("%s/monkeycode/vsix_version.txt", e.config.Extension.Baseurl)
	logger.With("url", u).Debug("开始获取远程插件版本信息")
	resp, err := http.Get(u)
	if err != nil {
		logger.With("error", err).Error("获取插件版本信息失败")
		return "", err
	}
	defer resp.Body.Close()

	version, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.With("error", err).Error("读取插件body信息失败")
		return "", err
	}

	logger.With("version", string(version)).Debug("读取到的插件版本信息")
	return string(version), nil
}
