TRUNCATE TABLE `promotion_template_models`;
TRUNCATE TABLE `user_coupon_models`;
INSERT INTO promotion_template_models (template_group_id, version, name, description, promotion_type, rule_definition, discount_type, discount_properties, start_date, end_date, is_exclusive, priority, is_active, created_at, updated_at)
VALUES
    ('group-new-user-100-20', 1, '新用户专享券', '新注册用户可领取的满100减20元优惠券', 'PLATFORM_SALE', 'fact.User.Labels.exists(label, label == "new_user")', 'FIXED_AMOUNT', '{"threshold": 10000, "amount": 2000}', '2025-01-01 00:00:00', '2025-12-31 23:59:59', 1, 100, 1, NOW(), NOW());
INSERT INTO promotion_template_models (template_group_id, version, name, description, promotion_type, rule_definition, discount_type, discount_properties, start_date, end_date, is_exclusive, priority, is_active, created_at, updated_at)
VALUES
    ('group-vip-88-percent', 1, 'VIP会员88折券', 'VIP会员专享，无门槛88折，最高可优惠50元', 'STORE_COUPON', 'fact.User.IsVip == true', 'PERCENTAGE', '{"percentage": 88, "ceiling": 5000}', '2025-01-01 00:00:00', '2025-12-31 23:59:59', 0, 90, 1, NOW(), NOW());
INSERT INTO promotion_template_models (template_group_id, version, name, description, promotion_type, rule_definition, discount_type, discount_properties, start_date, end_date, is_exclusive, priority, is_active, created_at, updated_at)
VALUES
    ('group-general-5', 1, '全场通用券', '已失效的无门槛5元券', 'PLATFORM_SALE', '', 'FIXED_AMOUNT', '{"threshold": 0, "amount": 500}', '2024-01-01 00:00:00', '2024-12-31 23:59:59', 0, 10, 0, NOW(), NOW());
INSERT INTO user_coupon_models (user_id, coupon_code, template_id, status, issue_date, expiry_date, created_at, updated_at)
VALUES
    (123, 'VIP-COUPON-123', 2, 'UNUSED', NOW(), '2025-12-31 23:59:59', NOW(), NOW());
INSERT INTO user_coupon_models (user_id, coupon_code, template_id, status, issue_date, expiry_date, created_at, updated_at)
VALUES
    (456, 'NEW-USER-COUPON-456', 1, 'UNUSED', NOW(), '2025-12-31 23:59:59', NOW(), NOW());