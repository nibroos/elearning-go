INSERT INTO
  pools (
    group1_id,
    group2_id,
    mv1_id,
    mv2_id,
    created_by_id,
    updated_by_id,
    created_at,
    updated_at
  )
SELECT
  (
    SELECT
      id
    FROM
      groups
    WHERE
      name = 'roles'
  ),
  (
    SELECT
      id
    FROM
      groups
    WHERE
      name = 'permissions'
  ),
  (
    SELECT
      id
    FROM
      mix_values
    WHERE
      name = 'superadmin'
  ),
  id,
  1,
  1,
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM
  mix_values
WHERE
  group_id = (
    SELECT
      id
    FROM
      groups
    WHERE
      name = 'permissions'
  );

COMMIT;