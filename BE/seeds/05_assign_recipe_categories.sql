-- Assign dish_type_id, region_id, main_ingredient_id, and meal_type_id
-- to existing recipes that were synced from Notion without category properties.
--
-- This maps recipes to categories by name patterns.
-- Future Notion syncs should set the Category/Region/Main Ingredient/Meal Type
-- properties directly in Notion, so this seed is a one-time data fix.

-- Phở Bò → Phở (DISH_TYPE), Bò (MAIN_INGREDIENT), Miền Bắc (REGION), Bữa tối (MEAL_TYPE)
UPDATE recipes SET
  dish_type_id       = (SELECT id FROM categories WHERE slug = 'pho'     AND type = 'DISH_TYPE'),
  main_ingredient_id = (SELECT id FROM categories WHERE slug = 'bo'      AND type = 'MAIN_INGREDIENT'),
  region_id          = (SELECT id FROM categories WHERE slug = 'mien-bac' AND type = 'REGION'),
  meal_type_id       = (SELECT id FROM categories WHERE slug = 'bua-toi' AND type = 'MEAL_TYPE')
WHERE slug = 'pho-bo' AND dish_type_id IS NULL;

-- Bún Chả → Bún (DISH_TYPE), Heo (MAIN_INGREDIENT), Miền Bắc (REGION), Bữa trưa (MEAL_TYPE)
UPDATE recipes SET
  dish_type_id       = (SELECT id FROM categories WHERE slug = 'bun'      AND type = 'DISH_TYPE'),
  main_ingredient_id = (SELECT id FROM categories WHERE slug = 'heo'      AND type = 'MAIN_INGREDIENT'),
  region_id          = (SELECT id FROM categories WHERE slug = 'mien-bac' AND type = 'REGION'),
  meal_type_id       = (SELECT id FROM categories WHERE slug = 'bua-trua' AND type = 'MEAL_TYPE')
WHERE slug = 'bun-cha' AND dish_type_id IS NULL;

-- Cơm Tấm → Cơm (DISH_TYPE), Heo (MAIN_INGREDIENT), Miền Nam (REGION), Bữa trưa (MEAL_TYPE)
UPDATE recipes SET
  dish_type_id       = (SELECT id FROM categories WHERE slug = 'com'      AND type = 'DISH_TYPE'),
  main_ingredient_id = (SELECT id FROM categories WHERE slug = 'heo'      AND type = 'MAIN_INGREDIENT'),
  region_id          = (SELECT id FROM categories WHERE slug = 'mien-nam' AND type = 'REGION'),
  meal_type_id       = (SELECT id FROM categories WHERE slug = 'bua-trua' AND type = 'MEAL_TYPE')
WHERE slug = 'com-tam' AND dish_type_id IS NULL;

-- Bánh Mì → Bánh (DISH_TYPE), Heo (MAIN_INGREDIENT), Miền Nam (REGION), Bữa sáng (MEAL_TYPE)
UPDATE recipes SET
  dish_type_id       = (SELECT id FROM categories WHERE slug = 'banh'     AND type = 'DISH_TYPE'),
  main_ingredient_id = (SELECT id FROM categories WHERE slug = 'heo'      AND type = 'MAIN_INGREDIENT'),
  region_id          = (SELECT id FROM categories WHERE slug = 'mien-nam' AND type = 'REGION'),
  meal_type_id       = (SELECT id FROM categories WHERE slug = 'bua-sang' AND type = 'MEAL_TYPE')
WHERE slug = 'banh-mi' AND dish_type_id IS NULL;

-- Gỏi Cuốn → Gỏi (DISH_TYPE), Tôm (MAIN_INGREDIENT), Miền Nam (REGION), Ăn vặt (MEAL_TYPE)
UPDATE recipes SET
  dish_type_id       = (SELECT id FROM categories WHERE slug = 'goi'      AND type = 'DISH_TYPE'),
  main_ingredient_id = (SELECT id FROM categories WHERE slug = 'tom'      AND type = 'MAIN_INGREDIENT'),
  region_id          = (SELECT id FROM categories WHERE slug = 'mien-nam' AND type = 'REGION'),
  meal_type_id       = (SELECT id FROM categories WHERE slug = 'an-vat'   AND type = 'MEAL_TYPE')
WHERE slug = 'goi-cuon' AND dish_type_id IS NULL;
