package domain

import "context"

// UnitOfWork 定义了工作单元的接口
// 它提供了一种方式来确保多个仓储操作在同一个事务中执行
type UnitOfWork interface {
	// Execute 将一个函数包裹在单个事务中执行
	// fn 是包含所有业务逻辑和仓储操作的函数
	// 如果 fn 返回错误，事务将回滚；否则，事务将提交
	Execute(ctx context.Context, fn func(repoProvider RepositoryProvider) error) error
}

// RepositoryProvider 是一个接口，用于在事务中获取仓储实例
// 这样可以确保所有获取到的仓储都共享同一个事务
type RepositoryProvider interface {
	Coupons() CouponRepository
	Templates() PromotionTemplateRepository
}
