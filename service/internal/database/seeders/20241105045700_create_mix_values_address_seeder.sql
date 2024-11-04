BEGIN;

INSERT INTO
  groups (
    id,
    name,
    description,
    status,
    created_by_id,
    updated_by_id,
    created_at,
    updated_at
  )
VALUES
  (
    5,
    'addresses',
    'Alamat',
    1,
    1,
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

INSERT INTO
  mix_values (
    group_id,
    name,
    description,
    status,
    created_at,
    updated_at
  )
VALUES
  (
    5,
    'Rumah',
    'Alamat - Rumah',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    5,
    'Kantor',
    'Alamat - Kantor',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    5,
    'Saudara',
    'Alamat - Saudara',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    5,
    'Lainnya',
    'Alamat - Lainnya',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;