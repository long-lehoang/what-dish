INSERT INTO tags (name, slug) VALUES
    ('Nhanh gọn', 'nhanh-gon'),
    ('Healthy', 'healthy'),
    ('Ít calo', 'it-calo'),
    ('Nhiều protein', 'nhieu-protein'),
    ('Cho trẻ em', 'cho-tre-em'),
    ('Tiệc', 'tiec'),
    ('Giảm cân', 'giam-can'),
    ('Tăng cơ', 'tang-co'),
    ('Chay', 'chay'),
    ('Không cay', 'khong-cay'),
    ('Cay', 'cay'),
    ('Dễ làm', 'de-lam'),
    ('Gia đình', 'gia-dinh'),
    ('Một mình', 'mot-minh'),
    ('Tiết kiệm', 'tiet-kiem')
ON CONFLICT (slug) DO NOTHING;
