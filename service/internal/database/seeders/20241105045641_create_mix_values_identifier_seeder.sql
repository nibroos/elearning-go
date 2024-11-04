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
    4,
    'identifiers',
    'Identifier (ID)',
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
    4,
    'KTP',
    'ID - Kartu Tanda Penduduk',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    4,
    'NPWP',
    'ID - Nomor Pokok Wajib Pajak',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    4,
    'SIM',
    'ID - Surat Izin Mengeemudi',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    4,
    'Passport',
    'ID - Passport',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    4,
    'Lainnya',
    'ID - Lainnya',
    1,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;