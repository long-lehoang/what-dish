-- Dish types
INSERT INTO categories (name, slug, type, sort_order) VALUES
    ('Kho', 'kho', 'DISH_TYPE', 1),
    ('Xào', 'xao', 'DISH_TYPE', 2),
    ('Hấp', 'hap', 'DISH_TYPE', 3),
    ('Nướng', 'nuong', 'DISH_TYPE', 4),
    ('Chiên', 'chien', 'DISH_TYPE', 5),
    ('Luộc', 'luoc', 'DISH_TYPE', 6),
    ('Canh', 'canh', 'DISH_TYPE', 7),
    ('Lẩu', 'lau', 'DISH_TYPE', 8),
    ('Gỏi', 'goi', 'DISH_TYPE', 9),
    ('Nộm', 'nom', 'DISH_TYPE', 10),
    ('Cháo', 'chao', 'DISH_TYPE', 11),
    ('Phở', 'pho', 'DISH_TYPE', 12),
    ('Bún', 'bun', 'DISH_TYPE', 13),
    ('Cơm', 'com', 'DISH_TYPE', 14),
    ('Bánh', 'banh', 'DISH_TYPE', 15)
ON CONFLICT (slug) DO NOTHING;

-- Regions
INSERT INTO categories (name, slug, type, sort_order) VALUES
    ('Miền Bắc', 'mien-bac', 'REGION', 1),
    ('Miền Trung', 'mien-trung', 'REGION', 2),
    ('Miền Nam', 'mien-nam', 'REGION', 3)
ON CONFLICT (slug) DO NOTHING;

-- Main ingredients
INSERT INTO categories (name, slug, type, sort_order) VALUES
    ('Gà', 'ga', 'MAIN_INGREDIENT', 1),
    ('Bò', 'bo', 'MAIN_INGREDIENT', 2),
    ('Heo', 'heo', 'MAIN_INGREDIENT', 3),
    ('Hải sản', 'hai-san', 'MAIN_INGREDIENT', 4),
    ('Cá', 'ca', 'MAIN_INGREDIENT', 5),
    ('Tôm', 'tom', 'MAIN_INGREDIENT', 6),
    ('Đậu phụ', 'dau-phu', 'MAIN_INGREDIENT', 7),
    ('Rau củ', 'rau-cu', 'MAIN_INGREDIENT', 8),
    ('Trứng', 'trung', 'MAIN_INGREDIENT', 9)
ON CONFLICT (slug) DO NOTHING;

-- Meal types
INSERT INTO categories (name, slug, type, sort_order) VALUES
    ('Bữa sáng', 'bua-sang', 'MEAL_TYPE', 1),
    ('Bữa trưa', 'bua-trua', 'MEAL_TYPE', 2),
    ('Bữa tối', 'bua-toi', 'MEAL_TYPE', 3),
    ('Ăn vặt', 'an-vat', 'MEAL_TYPE', 4)
ON CONFLICT (slug) DO NOTHING;
