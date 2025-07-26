package infrastructure

import (
	"context"
	"fmt"
	"github.com/wangyingjie930/nexus-promotion/internal/domain"
	"gorm.io/gorm"
)

// gormUnitOfWork 是 UnitOfWork 接口的GORM实现
type gormUnitOfWork struct {
	db *gorm.DB
}

// NewGormUnitOfWork 创建一个新的 GORM 工作单元实例
func NewGormUnitOfWork(db *gorm.DB) domain.UnitOfWork {
	return &gormUnitOfWork{db: db}
}

// Execute 实现了 domain.UnitOfWork 接口
func (uow *gormUnitOfWork) Execute(ctx context.Context, fn func(rp domain.RepositoryProvider) error) error {
	// 开始一个GORM事务
	tx := uow.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 创建一个使用该事务的仓储提供者
	repoProvider := &gormRepoProvider{db: tx}

	// 执行业务逻辑
	err := fn(repoProvider)
	if err != nil {
		// 如果业务逻辑返回错误，回滚事务
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rbErr, err)
		}
		return err // 返回原始的业务错误
	}

	// 如果业务逻辑成功，提交事务
	if cmtErr := tx.Commit().Error; cmtErr != nil {
		return fmt.Errorf("commit error: %w", cmtErr)
	}

	return nil
}

// gormRepoProvider 实现了 RepositoryProvider 接口
type gormRepoProvider struct {
	db *gorm.DB
}

func (p *gormRepoProvider) Coupons() domain.CouponRepository {
	// 注意：这里传入的是事务句柄 tx
	return NewGormCouponRepository(p.db)
}

func (p *gormRepoProvider) Templates() domain.PromotionTemplateRepository {
	// 注意：这里传入的是事务句柄 tx
	return NewGormPromotionTemplateRepository(p.db)
}
