-- +goose Up
ALTER TABLE product.product         ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);
ALTER TABLE product.product_variant ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);
ALTER TABLE product.product_collection ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);
ALTER TABLE product.product_category   ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);

-- +goose Down
ALTER TABLE product.product_category   DROP COLUMN IF EXISTS last_updated_by;
ALTER TABLE product.product_collection DROP COLUMN IF EXISTS last_updated_by;
ALTER TABLE product.product_variant    DROP COLUMN IF EXISTS last_updated_by;
ALTER TABLE product.product            DROP COLUMN IF EXISTS last_updated_by;
