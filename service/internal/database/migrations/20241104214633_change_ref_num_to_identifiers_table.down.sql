BEGIN;

ALTER TABLE
  addresses
ALTER COLUMN
  ref_num TYPE INT USING ref_num :: INTEGER;

ALTER TABLE
  identifiers
ALTER COLUMN
  ref_num TYPE INT USING ref_num :: INTEGER;

ALTER TABLE
  contacts
ALTER COLUMN
  ref_num TYPE INT USING ref_num :: INTEGER;

COMMIT;