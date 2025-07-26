// promotion-service/internal/domain/fact.go
package domain

import "time"

// CartItem 代表购物车中的一个商品项
type CartItem struct {
	SKU      string `json:"SKU"`
	Price    int64  `json:"Price"` // 商品单价（单位：分）
	Quantity int32  `json:"Quantity"`
	Category string `json:"Category"` // 商品品类
	Brand    string `json:"Brand"`    // 商品品牌
}

// UserContext 代表当前的用户信息
type UserContext struct {
	ID     int64    `json:"ID"`
	IsVip  bool     `json:"IsVip"`  // 是否为VIP用户
	Labels []string `json:"Labels"` // 用户标签，如 "new_user", "high_value"
}

// EnvironmentContext 代表环境信息
type EnvironmentContext struct {
	Timestamp time.Time `json:"Timestamp"` // 当前时间
	Channel   string    `json:"Channel"`   // 渠道, e.g., "app", "mini_program"
}

// Fact 是规则引擎和优惠计算策略所需的所有上下文信息的集合。
// 它是一个高度结构化的数据对象，作为评估过程的唯一输入。
// 这种设计将计算逻辑与数据来源完全解耦，极大地提高了系统的可测试性和可扩展性。
type Fact struct {
	User        UserContext        `json:"User"`
	Items       []CartItem         `json:"Items"`
	Environment EnvironmentContext `json:"Environment"`

	// 派生字段，在服务层预先计算，以简化规则逻辑
	TotalAmount int64 `json:"TotalAmount"` // 购物车总金额
}

// Trait 常量定义，用于 CEL 规则引擎的类型检查
const (
	TraitComparable = 1 // 可比较的
	TraitIterable   = 2 // 可迭代的
	TraitIndexable  = 3 // 可索引的
	TraitCallable   = 4 // 可调用的
)

// HasTrait returns whether the type has a given trait associated with it.
//
// 这个方法用于 CEL 规则引擎的类型检查，支持以下 trait：
// - TraitComparable: 支持比较操作
// - TraitIterable: 支持迭代操作
// - TraitIndexable: 支持索引操作
// - TraitCallable: 支持函数调用
func (f Fact) HasTrait(trait int) bool {
	switch trait {
	case TraitComparable:
		return true // Fact 类型支持比较操作
	case TraitIterable:
		return false // Fact 类型不支持迭代操作
	case TraitIndexable:
		return false // Fact 类型不支持索引操作
	case TraitCallable:
		return false // Fact 类型不支持函数调用
	default:
		return false // 未知的 trait 返回 false
	}
}

// TypeName returns the qualified type name of the type.
//
// The type name is also used as the type's identifier name at type-check and interpretation time.
// 这个方法用于 CEL 规则引擎的类型识别和错误报告。
func (f Fact) TypeName() string {
	return "domain.Fact"
}
