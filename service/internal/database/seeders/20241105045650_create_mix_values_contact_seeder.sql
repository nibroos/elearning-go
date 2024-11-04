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
    6,
    'contacts',
    'Kontak',
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
    6,
    'Pribadi (Nomor HP/WhatsApp)',
    'Kontak - Pribadi',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    6,
    'Rumah',
    'Kontak - Rumah',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    6,
    'Kantor',
    'Kontak - Kantor',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    6,
    'Saudara',
    'Kontak - Saudara',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;