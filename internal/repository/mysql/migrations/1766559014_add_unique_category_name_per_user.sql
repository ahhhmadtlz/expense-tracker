-- +migrate Up
-- First, remove duplicate categories (keep only the oldest one for each user+name)
DELETE c1 FROM categories c1
INNER JOIN categories c2 
WHERE c1.id > c2.id 
  AND c1.user_id = c2.user_id 
  AND c1.name = c2.name;

-- Then add the unique constraint
ALTER TABLE `categories` 
ADD UNIQUE KEY `unique_user_category_name` (`user_id`, `name`);

-- +migrate Down
ALTER TABLE `categories` 
DROP INDEX `unique_user_category_name`;
