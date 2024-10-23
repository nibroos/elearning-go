-- Insert user-role relationships into pools table
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
VALUES
  (
    (
      SELECT
        id
      FROM
        groups
      WHERE
        name = 'users'
    ),
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
        users
      WHERE
        username = 'nibros'
    ),
    (
      SELECT
        id
      FROM
        mix_values
      WHERE
        name = 'superadmin'
    ),
    1,
    1,
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
        name = 'users'
    ),
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
        users
      WHERE
        username = 'user1'
    ),
    (
      SELECT
        id
      FROM
        mix_values
      WHERE
        name = 'student'
    ),
    1,
    1,
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
        name = 'users'
    ),
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
        users
      WHERE
        username = 'user2'
    ),
    (
      SELECT
        id
      FROM
        mix_values
      WHERE
        name = 'lecturer'
    ),
    1,
    1,
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
        name = 'users'
    ),
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
        users
      WHERE
        username = 'user3'
    ),
    (
      SELECT
        id
      FROM
        mix_values
      WHERE
        name = 'student'
    ),
    1,
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;