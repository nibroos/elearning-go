BEGIN;

INSERT INTO
  groups (
    name,
    description,
    status,
    options_json,
    created_by_id,
    updated_by_id
  )
VALUES
  (
    'roles',
    'Roles table for storing user data',
    1,
    '{}',
    1,
    1
  ),
  (
    'permissions',
    'Permissions table for storing user data',
    1,
    '{}',
    1,
    1
  ),
  (
    'users',
    'Users table for storing user data',
    1,
    '{}',
    1,
    1
  );

COMMIT;