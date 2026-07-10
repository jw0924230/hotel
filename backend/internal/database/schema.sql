-- Database Schema for Hotel Management System
-- You can run this in your PostgreSQL database (e.g. Supabase SQL Editor) to initialize the tables.

CREATE TABLE IF NOT EXISTS hotels (
    id VARCHAR(255) PRIMARY KEY,
    name TEXT NOT NULL,
    area TEXT,
    address TEXT,
    phone TEXT,
    fax TEXT,
    website TEXT,
    email TEXT,
    category TEXT,
    stay_info TEXT,
    housing_rules TEXT,
    description TEXT,
    booking_link TEXT,
    is_disabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for searching hotels by name and area
CREATE INDEX IF NOT EXISTS idx_hotels_name ON hotels(name);
CREATE INDEX IF NOT EXISTS idx_hotels_area ON hotels(area);

-- Users table for admin backend auth
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert preset admin account if it doesn't exist
-- Email: jw0924230@gmail.com
-- Password: j1082087J
INSERT INTO users (email, password, role)
VALUES ('jw0924230@gmail.com', '$2a$10$suLYHfhAZSyncnjarD/bHuiROr1Fzu9K96HMND/nUaxsiCKigMAte', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Enlarge columns in case table already exists with small limits
ALTER TABLE hotels ALTER COLUMN id TYPE VARCHAR(255);
ALTER TABLE hotels ALTER COLUMN name TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN area TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN address TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN phone TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN fax TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN website TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN email TYPE TEXT;
ALTER TABLE hotels ALTER COLUMN category TYPE TEXT;
ALTER TABLE hotels ADD COLUMN IF NOT EXISTS is_disabled BOOLEAN DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS hotel_prices (
    hotel_id VARCHAR(255) PRIMARY KEY REFERENCES hotels(id) ON DELETE CASCADE,
    weekday_stay INT,
    holiday_stay INT,
    weekday_rest_hours INT,
    weekday_rest INT,
    holiday_rest_hours INT,
    holiday_rest INT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (weekday_stay IS NULL OR weekday_stay >= 0),
    CHECK (holiday_stay IS NULL OR holiday_stay >= 0),
    CHECK (weekday_rest_hours IS NULL OR weekday_rest_hours >= 0),
    CHECK (weekday_rest IS NULL OR weekday_rest >= 0),
    CHECK (holiday_rest_hours IS NULL OR holiday_rest_hours >= 0),
    CHECK (holiday_rest IS NULL OR holiday_rest >= 0)
);

ALTER TABLE hotel_prices ADD COLUMN IF NOT EXISTS weekday_rest_hours INT;
ALTER TABLE hotel_prices ADD COLUMN IF NOT EXISTS holiday_rest_hours INT;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'hotel_prices' AND column_name = 'rest_hours'
    ) THEN
        EXECUTE $migration$
            UPDATE hotel_prices
            SET weekday_rest_hours = COALESCE(weekday_rest_hours, ROUND(rest_hours)::INT),
                holiday_rest_hours = COALESCE(holiday_rest_hours, ROUND(rest_hours)::INT)
        $migration$;
    END IF;
END $$;

ALTER TABLE hotel_prices DROP COLUMN IF EXISTS rest_hours;

-- Extract legacy free-form price strings into structured price columns.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'hotels' AND column_name = 'price_accommodation'
    ) THEN
        EXECUTE $migration$
            INSERT INTO hotel_prices (
                hotel_id, weekday_stay, holiday_stay, weekday_rest_hours,
                weekday_rest, holiday_rest_hours, holiday_rest
            )
            SELECT
                h.id,
                CASE WHEN cardinality(stay_numbers.values) >= 1
                    THEN ROUND(stay_numbers.values[1])::INT END,
                CASE WHEN cardinality(stay_numbers.values) >= 2
                    THEN ROUND(stay_numbers.values[2])::INT END,
                CASE WHEN cardinality(rest_numbers.values) >= 3
                    THEN ROUND(rest_numbers.values[1])::INT END,
                CASE
                    WHEN cardinality(rest_numbers.values) >= 3
                        THEN ROUND(rest_numbers.values[2])::INT
                    WHEN cardinality(rest_numbers.values) >= 1
                        THEN ROUND(rest_numbers.values[1])::INT
                END,
                CASE WHEN cardinality(rest_numbers.values) >= 3
                    THEN ROUND(rest_numbers.values[1])::INT END,
                CASE
                    WHEN h.price_rest LIKE '%假日%'
                         AND cardinality(rest_numbers.values) >= 3
                        THEN ROUND(rest_numbers.values[3])::INT
                    WHEN h.price_rest LIKE '%假日%'
                         AND cardinality(rest_numbers.values) >= 2
                        THEN ROUND(rest_numbers.values[2])::INT
                END
            FROM hotels h
            CROSS JOIN LATERAL (
                SELECT ARRAY(
                    SELECT match[1]::NUMERIC
                    FROM regexp_matches(
                        COALESCE(h.price_accommodation, ''),
                        '([0-9]+(?:\.[0-9]+)?)',
                        'g'
                    ) AS match
                ) AS values
            ) stay_numbers
            CROSS JOIN LATERAL (
                SELECT ARRAY(
                    SELECT match[1]::NUMERIC
                    FROM regexp_matches(
                        COALESCE(h.price_rest, ''),
                        '([0-9]+(?:\.[0-9]+)?)',
                        'g'
                    ) AS match
                ) AS values
            ) rest_numbers
            WHERE cardinality(stay_numbers.values) > 0
               OR cardinality(rest_numbers.values) > 0
            ON CONFLICT (hotel_id) DO NOTHING
        $migration$;
    END IF;
END $$;

ALTER TABLE hotels DROP COLUMN IF EXISTS price_accommodation;
ALTER TABLE hotels DROP COLUMN IF EXISTS price_rest;
ALTER TABLE hotels DROP COLUMN IF EXISTS price;

-- Hotel images are normalized so ordering and image metadata can be managed independently.
CREATE TABLE IF NOT EXISTS hotel_images (
    id BIGSERIAL PRIMARY KEY,
    hotel_id VARCHAR(255) NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE hotel_images DROP CONSTRAINT IF EXISTS hotel_images_hotel_id_url_key;

CREATE INDEX IF NOT EXISTS idx_hotel_images_hotel_sort
ON hotel_images(hotel_id, sort_order, id);

-- Migrate legacy JSONB image arrays before removing the old column.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'hotels' AND column_name = 'images'
    ) THEN
        EXECUTE $migration$
            INSERT INTO hotel_images (hotel_id, url, sort_order)
            SELECT DISTINCT ON (h.id, image.url)
                h.id, image.url, image.ordinality - 1
            FROM hotels h
            CROSS JOIN LATERAL jsonb_array_elements_text(COALESCE(h.images, '[]'::jsonb))
                WITH ORDINALITY AS image(url, ordinality)
            WHERE image.url <> ''
            ORDER BY h.id, image.url, image.ordinality
            ON CONFLICT (hotel_id, url) DO UPDATE
            SET sort_order = EXCLUDED.sort_order
        $migration$;
    END IF;
END $$;

ALTER TABLE hotels DROP COLUMN IF EXISTS images;
ALTER TABLE hotels DROP COLUMN IF EXISTS precautions;

-- General-purpose categories table (type-based: city, hotel_category, etc.)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,        -- 'city', 'hotel_category', etc.
    name TEXT NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(type, name)
);

CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(type);

-- Remove the legacy city tables after the unified categories table exists.
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS cites;

-- Seed city data (sort_order preserves original city IDs for backward compat)
INSERT INTO categories (type, name, sort_order) VALUES
('city', '基隆', 1), ('city', '台北', 2), ('city', '新北', 3),
('city', '桃園', 4), ('city', '新竹', 5), ('city', '宜蘭', 6),
('city', '苗栗', 7), ('city', '台中', 8), ('city', '彰化', 10),
('city', '南投', 11), ('city', '雲林', 12), ('city', '嘉義', 13),
('city', '台南', 14), ('city', '高雄', 16), ('city', '屏東', 18),
('city', '花蓮', 19), ('city', '台東', 20), ('city', '澎湖', 21),
('city', '金門', 22), ('city', '馬祖', 23), ('city', '其他', 24)
ON CONFLICT (type, name) DO UPDATE SET sort_order = EXCLUDED.sort_order;

-- Seed hotel category data
DELETE FROM categories
WHERE type = 'hotel_category'
  AND name NOT IN ('汽車旅館', '精品商旅', '溫泉會館');

INSERT INTO categories (type, name, sort_order) VALUES
('hotel_category', '汽車旅館', 1),
('hotel_category', '精品商旅', 2),
('hotel_category', '溫泉會館', 3)
ON CONFLICT (type, name) DO UPDATE SET sort_order = EXCLUDED.sort_order;

-- Blog Posts table for article publishing and SEO management
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    tags JSONB DEFAULT '[]'::jsonb, -- Array of tags, e.g., ["住宿推薦", "大安區"]
    image TEXT,                    -- Main cover image URL (Imgur or local)
    content TEXT NOT NULL,         -- Rich text HTML from wangEditor
    ad_link TEXT,                  -- Optional HTML ad links
    seo_title TEXT,                -- SEO Page Title
    seo_keywords TEXT,             -- SEO Keywords
    seo_description TEXT,          -- SEO Description
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for searching posts
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts(title);

-- Homepage featured hotels slots selection
CREATE TABLE IF NOT EXISTS homepage_hotels (
    city VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL,
    hotel_id VARCHAR(255) NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    PRIMARY KEY (city, sort_order)
);
