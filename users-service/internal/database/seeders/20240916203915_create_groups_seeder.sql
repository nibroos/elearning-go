BEGIN;

INSERT INTO
  groups (
    id,
    name,
    description,
    status,
    options_json,
    created_by_id,
    updated_by_id
  )
VALUES
  (
    1,
    'roles',
    'Roles table for storing user data',
    1,
    '{}',
    1,
    1
  ),
  (
    2,
    'permissions',
    'Permissions table for storing user data',
    1,
    '{}',
    1,
    1
  ),
  (
    3,
    'users',
    'Users table for storing user data',
    1,
    '{}',
    1,
    1
  );

COMMIT;