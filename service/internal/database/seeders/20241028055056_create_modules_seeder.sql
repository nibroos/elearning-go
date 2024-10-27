BEGIN;

-- ID uint `json:"id" gorm:"column:id;primaryKey;autoIncrement"` Name string `json:"name" gorm:"column:name"` Description string `json:"desc" gorm:"column:description"` LogoURL string `json:"logo_url" gorm:"column:logo_url"` VideoURL string `json:"video_url" gorm:"column:video_url"` CreatedByID * uint `json:"created_by_id" gorm:"column:created_by_id"` UpdatedByID * uint `json:"updated_by_id" gorm:"column:updated_by_id"`
INSERT INTO
  modules (
    id,
    class_id,
    name,
    description,
    created_by_id,
    updated_by_id
  )
VALUES
  (
    1,
    1,
    'A',
    'A modules',
    1,
    1
  ),
  (
    2,
    2,
    'B',
    'B modules',
    1,
    1
  ),
  (
    3,
    3,
    'C',
    'C modules',
    1,
    1
  );

-- Reset the sequence to the maximum id value
-- Replace `modules_id_seq` with the actual sequence name if different
SELECT
  setval(
    'modules_id_seq',
    (
      SELECT
        MAX(id)
      FROM
        modules
    )
  );

COMMIT;