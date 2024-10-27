BEGIN;

INSERT INTO
  classes (
    id,
    subscribe_id,
    incharge_id,
    name,
    description,
    banner_url,
    logo_url,
    video_url,
    created_by_id,
    updated_by_id
  )
VALUES
  (
    1,
    1,
    1,
    'A Classes',
    'A classes description',
    '',
    '',
    '',
    1,
    1
  ),
  (
    2,
    2,
    1,
    'B Classes',
    'B classes description',
    '',
    '',
    '',
    1,
    1
  ),
  (
    3,
    3,
    1,
    'C Classes',
    'C classes description',
    '',
    '',
    '',
    1,
    1
  );

-- Reset the sequence to the maximum id value
-- Replace `subscribes_id_seq` with the actual sequence name if different
SELECT
  setval(
    'classes_id_seq',
    (
      SELECT
        MAX(id)
      FROM
        classes
    )
  );

COMMIT;