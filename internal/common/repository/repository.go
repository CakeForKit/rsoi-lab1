package repository

import (
	"context"
	"fmt"
	"strings"
	"uuid"

	"github.com/CakeForKit/rsoi-lab1/internal/common/db"
	"github.com/CakeForKit/rsoi-lab1/internal/common/utils"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	GetById(ctx context.Context, ids []uuid.UUID) ([]T, error)
	GetAll(ctx context.Context) ([]T, error)

	Create(ctx context.Context, entities []T) ([]T, error)
	Update(ctx context.Context, entities []T) ([]T, error)
	DeleteById(ctx context.Context, ids []uuid.UUID) error
}

type postgresRepository[T any] struct {
	dataSource *gorm.DB
}

func NewPostgresRepository[T any]() (Repository[T], error) {
	dataSource := db.GetPostgresDataSource()
	var entity T
	if err := dataSource.AutoMigrate(&entity); err != nil && !strings.Contains(err.Error(), "already exists") {
		return nil, err
	}
	return &postgresRepository[T]{dataSource: dataSource}, nil
}

func (pgRepo *postgresRepository[T]) GetById(ctx context.Context, ids []uuid.UUID) ([]T, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids can't be empty")
	}

	var result []T
	tx := pgRepo.dataSource.WithContext(ctx)
	if err := tx.Where("id in ?", utils.Map(ids, func(id uuid.UUID) string { return id.String() })).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (pgRepo *postgresRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var result []T
	tx := pgRepo.dataSource.WithContext(ctx)
	if err := tx.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (pgRepo *postgresRepository[T]) Create(ctx context.Context, entities []T) ([]T, error) {
	tx := pgRepo.dataSource.WithContext(ctx)
	result := tx.Create(&entities)
	if err := result.Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (pgRepo *postgresRepository[T]) Update(ctx context.Context, entities []T) ([]T, error) {
	tx := pgRepo.dataSource.WithContext(ctx)
	result := tx.Save(&entities)
	if err := result.Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (pgRepo *postgresRepository[T]) DeleteById(ctx context.Context, ids []uuid.UUID) error {
	var entity []T
	tx := pgRepo.dataSource.WithContext(ctx)
	result := tx.Where("id in ?", utils.Map(ids, func(it uuid.UUID) string { return it.String() })).Delete(&entity)
	return result.Error
}
