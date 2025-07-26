// internal/interfaces/http_handler_test.go
package interfaces

import (
	"bytes"
	"encoding/json"
	"github.com/wangyingjie930/nexus-promotion/internal/application"
	"github.com/wangyingjie930/nexus-promotion/internal/domain"
	"github.com/wangyingjie930/nexus-promotion/internal/infrastructure"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setupTestServer 初始化一个完整的测试服务器，包含所有依赖项
func setupTestServer(t *testing.T) *httptest.Server {
	// 使用测试专用的数据库或配置
	// 注意：为了安全，通常从环境变量或配置文件中读取测试数据库连接信息
	dsn := "root:root@tcp(mysql.infra:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// 自动迁移数据库模型
	err = db.AutoMigrate(&infrastructure.PromotionTemplateModel{}, &infrastructure.UserCouponModel{})
	if err != nil {
		t.Fatalf("failed to auto migrate models: %v", err)
	}

	// 依赖注入，与main.go中的逻辑保持一致
	couponRepo := infrastructure.NewGormCouponRepository(db)
	templateRepo := infrastructure.NewGormPromotionTemplateRepository(db)
	uow := infrastructure.NewGormUnitOfWork(db)
	tracer := otel.Tracer("test-tracer")
	promoService := application.NewPromotionService(uow, templateRepo, couponRepo, tracer)
	promoHandler := NewPromotionHandler(promoService)

	// 创建 Mux 并注册路由
	mux := http.NewServeMux()
	promoHandler.RegisterRoutes(mux)

	// 使用 httptest 创建一个测试服务器
	return httptest.NewServer(mux)
}

// TestCalculateBestOffer_VipUser_Success 是一个具体的功能测试用例
func TestCalculateBestOffer_VipUser_Success(t *testing.T) {
	// 1. 设置
	server := setupTestServer(t)
	defer server.Close() // 确保测试结束后关闭服务器

	// 2. 准备请求数据 (Fact)
	fact := domain.Fact{
		User: domain.UserContext{
			ID:    123, // 对应 test_seed.sql 中的VIP用户
			IsVip: true,
		},
		Items: []domain.CartItem{
			{SKU: "SKU001", Price: 15000, Quantity: 1, Category: "Electronics"},
		},
		Environment: domain.EnvironmentContext{
			Timestamp: time.Now(),
			Channel:   "app",
		},
		TotalAmount: 15000, // 购物车总金额150元
	}
	body, _ := json.Marshal(fact)

	// 3. 发起HTTP请求
	resp, err := http.Post(server.URL+"/offers/calculate-best", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to send request to test server: %v", err)
	}
	defer resp.Body.Close()

	// 4. 断言 (Assert)
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	// 解析响应体
	var respBody application.DiscountApplicationResponse
	respBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// 核心业务断言：
	// 订单金额150元，VIP用户有88折券 (TemplateID=2)，无门槛
	// 优惠金额 = 15000 * (100 - 88) / 100 = 1800 (18元)
	expectedAmount := int64(1800)
	if respBody.Amount != expectedAmount {
		t.Errorf("expected discount amount %d; got %d", expectedAmount, respBody.Amount)
	}
	if respBody.StrategyName != "PercentageStrategy" {
		t.Errorf("expected strategy 'PercentageStrategy'; got '%s'", respBody.StrategyName)
	}
}
