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
    'permissions',
    'Permissions for the application',
    1,
    '{}',
    1,
    1
  );

COMMIT;