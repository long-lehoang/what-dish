INSERT INTO nutrition_goals (name, description, meal_calories_min, meal_calories_max, daily_calories_min, daily_calories_max, protein_pct, carbs_pct, fat_pct, sort_order) VALUES
    ('Giảm cân', 'Giảm cân lành mạnh, thâm hụt calo vừa phải', 300, 500, 1200, 1600, 30, 40, 30, 1),
    ('Duy trì', 'Duy trì cân nặng hiện tại', 400, 700, 1800, 2200, 25, 50, 25, 2),
    ('Tăng cân', 'Tăng cân và xây dựng cơ bắp', 500, 900, 2200, 3000, 30, 45, 25, 3),
    ('Eat Clean', 'Ăn sạch, thực phẩm tự nhiên', 350, 600, 1500, 2000, 30, 40, 30, 4)
ON CONFLICT DO NOTHING;
