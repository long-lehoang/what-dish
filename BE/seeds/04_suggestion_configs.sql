INSERT INTO suggestion_configs (config_type, name, description, params, sort_order) VALUES
    ('CALORIE_PRESET', 'Nhẹ nhàng', 'Bữa ăn nhẹ, ít calo', '{"min_cal": 200, "max_cal": 400}', 1),
    ('CALORIE_PRESET', 'Vừa phải', 'Bữa ăn tiêu chuẩn', '{"min_cal": 400, "max_cal": 700}', 2),
    ('CALORIE_PRESET', 'No bụng', 'Bữa ăn đầy đủ năng lượng', '{"min_cal": 700, "max_cal": 1000}', 3),
    ('GROUP_PRESET', 'Cặp đôi', 'Bữa ăn cho 2 người', '{"group_size": 2, "min_dishes": 2, "max_dishes": 3}', 1),
    ('GROUP_PRESET', 'Gia đình nhỏ', 'Bữa ăn cho 3-4 người', '{"group_size": 4, "min_dishes": 3, "max_dishes": 4}', 2),
    ('GROUP_PRESET', 'Gia đình lớn', 'Bữa ăn cho 5-8 người', '{"group_size": 6, "min_dishes": 4, "max_dishes": 5}', 3),
    ('GROUP_PRESET', 'Tiệc nhỏ', 'Bữa tiệc cho 8-12 người', '{"group_size": 10, "min_dishes": 5, "max_dishes": 7}', 4)
ON CONFLICT DO NOTHING;
