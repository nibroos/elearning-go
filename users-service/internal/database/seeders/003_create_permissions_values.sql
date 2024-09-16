BEGIN;

INSERT INTO
  mix_values (
    group_id,
    name,
    description,
    status,
    options_json,
    created_at,
    updated_at
  )
VALUES
  (
    (
      SELECT
        id
      FROM
        groups
      WHERE
        name = 'roles'
    ),
    'superadmin',
    'Super Admin Role',
    1,
    '{}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    (
      SELECT
        id
      FROM
        groups
      WHERE
        name = 'roles'
    ),
    'student',
    'Student Role',
    1,
    '{}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    (
      SELECT
        id
      FROM
        groups
      WHERE
        name = 'roles'
    ),
    'lecturer',
    'Lecturer Role',
    1,
    '{}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;