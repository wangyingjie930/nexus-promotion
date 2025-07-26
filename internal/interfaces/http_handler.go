package interfaces

import (
	"encoding/json"
	"github.com/wangyingjie930/nexus-pkg/logger"
	"github.com/wangyingjie930/nexus-promotion/internal/application"
	"github.com/wangyingjie930/nexus-promotion/internal/domain"
	"net/http"
	"strconv"
)

// PromotionHandler 封装了应用服务，并处理HTTP请求。
type PromotionHandler struct {
	promoService application.PromotionService
}

// NewPromotionHandler 创建一个新的PromotionHandler实例。
func NewPromotionHandler(promoService application.PromotionService) *PromotionHandler {
	return &PromotionHandler{
		promoService: promoService,
	}
}

// RegisterRoutes 将所有促销相关的HTTP路由注册到一个Mux上。
// 这里我们假设Mux是一个接口，拥有类似http.ServeMux的 HandleFunc 方法。
func (h *PromotionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /templates", h.CreatePromotionTemplate)
	mux.HandleFunc("PUT /templates", h.UpdatePromotionTemplate)
	mux.HandleFunc("DELETE /templates/{groupId}", h.DeactivatePromotionTemplate)
	mux.HandleFunc("GET /templates/{id}", h.GetPromotionTemplate)
	mux.HandleFunc("GET /templates/group/{groupId}", h.GetActiveTemplateByGroup)
	mux.HandleFunc("POST /coupons/issue", h.IssueCouponToUser)
	mux.HandleFunc("POST /coupons/issue-batch", h.IssueCouponsInBatch)
	mux.HandleFunc("POST /offers/calculate-best", h.CalculateBestOffer)
	mux.HandleFunc("POST /users/{userId}/applicable-coupons", h.GetApplicableCoupons)
	mux.HandleFunc("POST /users/{userId}/coupons/{couponCode}/freeze", h.FreezeUserCoupon)
	mux.HandleFunc("POST /users/{userId}/coupons/{couponCode}/use", h.UseUserCoupon)
	mux.HandleFunc("POST /users/{userId}/coupons/{couponCode}/unfreeze", h.UnfreezeUserCoupon)
}

// --- Handler 方法实现 ---

func (h *PromotionHandler) CreatePromotionTemplate(w http.ResponseWriter, r *http.Request) {
	var req application.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.promoService.CreatePromotionTemplate(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) UpdatePromotionTemplate(w http.ResponseWriter, r *http.Request) {
	var req application.UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.promoService.UpdatePromotionTemplate(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) DeactivatePromotionTemplate(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	if groupID == "" {
		http.Error(w, "groupId is required", http.StatusBadRequest)
		return
	}
	err := h.promoService.DeactivatePromotionTemplate(r.Context(), groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PromotionHandler) GetPromotionTemplate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid template ID", http.StatusBadRequest)
		return
	}
	resp, err := h.promoService.GetPromotionTemplate(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) GetActiveTemplateByGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.PathValue("groupId")
	if groupID == "" {
		http.Error(w, "groupId is required", http.StatusBadRequest)
		return
	}
	resp, err := h.promoService.GetActiveTemplateByGroup(r.Context(), groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) IssueCouponToUser(w http.ResponseWriter, r *http.Request) {
	var req application.IssueCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.promoService.IssueCouponToUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) IssueCouponsInBatch(w http.ResponseWriter, r *http.Request) {
	var req application.BatchIssueCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.promoService.IssueCouponsInBatch(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *PromotionHandler) CalculateBestOffer(w http.ResponseWriter, r *http.Request) {
	var fact domain.Fact
	if err := json.NewDecoder(r.Body).Decode(&fact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.promoService.CalculateBestOffer(r.Context(), &fact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var fact domain.Fact
	if err := json.NewDecoder(r.Body).Decode(&fact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Ctx(r.Context()).Info().Any("fact", fact).Send()

	resp, err := h.promoService.GetApplicableCoupons(r.Context(), &fact, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PromotionHandler) FreezeUserCoupon(w http.ResponseWriter, r *http.Request) {
	userID, couponCode := h.parseUserAndCouponParams(r)
	if userID == 0 || couponCode == "" {
		http.Error(w, "invalid user ID or coupon code", http.StatusBadRequest)
		return
	}
	err := h.promoService.FreezeUserCoupon(r.Context(), userID, couponCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PromotionHandler) UseUserCoupon(w http.ResponseWriter, r *http.Request) {
	userID, couponCode := h.parseUserAndCouponParams(r)
	if userID == 0 || couponCode == "" {
		http.Error(w, "invalid user ID or coupon code", http.StatusBadRequest)
		return
	}
	err := h.promoService.UseUserCoupon(r.Context(), userID, couponCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *PromotionHandler) UnfreezeUserCoupon(w http.ResponseWriter, r *http.Request) {
	userID, couponCode := h.parseUserAndCouponParams(r)
	if userID == 0 || couponCode == "" {
		http.Error(w, "invalid user ID or coupon code", http.StatusBadRequest)
		return
	}
	err := h.promoService.UnfreezeUserCoupon(r.Context(), userID, couponCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// parseUserAndCouponParams 是一个辅助函数，用于从URL路径中解析参数
func (h *PromotionHandler) parseUserAndCouponParams(r *http.Request) (int64, string) {
	userIDStr := r.PathValue("userId")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	couponCode := r.PathValue("couponCode")
	return userID, couponCode
}
