package bootstrap

import (
	"fmt"

	"github.com/carloscacb333/go-hexagonal/app/shared/infrastructure/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ContainerBuilder struct {
	config *config.Config
	logger *zap.Logger
	db     *gorm.DB
	errors []error
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		errors: make([]error, 0),
	}
}

func (b *ContainerBuilder) WithConfig(cfg *config.Config) *ContainerBuilder {
	if cfg == nil {
		b.errors = append(b.errors, fmt.Errorf("config cannot be nil"))
		return b
	}
	b.config = cfg
	return b
}

func (b *ContainerBuilder) WithLogger(logger *zap.Logger) *ContainerBuilder {
	if logger == nil {
		b.errors = append(b.errors, fmt.Errorf("logger cannot be nil"))
		return b
	}
	b.logger = logger
	return b
}

func (b *ContainerBuilder) WithDatabase(db *gorm.DB) *ContainerBuilder {
	if db == nil {
		b.errors = append(b.errors, fmt.Errorf("database cannot be nil"))
		return b
	}
	b.db = db
	return b
}

func (b *ContainerBuilder) Build() (*Container, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("container build failed with %d errors: %v", len(b.errors), b.errors)
	}

	if b.config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if b.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if b.db == nil {
		return nil, fmt.Errorf("database is required")
	}

	return NewContainer(b.config, b.logger, b.db)
}
