BEGIN;

INSERT INTO
  subscribes (
    id,
    name,
    description,
    created_by_id,
    updated_by_id
  )
VALUES
  (
    1,
    'A',
    'A subscribes',
    1,
    1
  ),
  (
    2,
    'B',
    'B subscribes',
    1,
    1
  ),
  (
    3,
    'C',
    'C subscribes',
    1,
    1
  );

-- Reset the sequence to the maximum id value
-- Replace `subscribes_id_seq` with the actual sequence name if different
SELECT
  setval(
    'subscribes_id_seq',
    (
      SELECT
        MAX(id)
      FROM
        subscribes
    )
  );

COMMIT;