BEGIN;

-- Insert users into users table
INSERT INTO
  users (
    username,
    email,
    name,
    password,
    address,
    created_at,
    updated_at
  )
VALUES
  (
    'nibros',
    'nibros@example.com',
    'Nibros',
    crypt('admel', gen_salt('bf')),
    '123 Main St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user1',
    'user1@example.com',
    'User One',
    crypt('password1', gen_salt('bf')),
    '456 Elm St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user2',
    'user2@example.com',
    'User Two',
    crypt('password2', gen_salt('bf')),
    '789 Oak St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user3',
    'user3@example.com',
    'User Three',
    crypt('password3', gen_salt('bf')),
    '101 Pine St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user4',
    'user4@example.com',
    'User Four',
    crypt('password4', gen_salt('bf')),
    '202 Maple St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user5',
    'user5@example.com',
    'User Five',
    crypt('password5', gen_salt('bf')),
    '303 Birch St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user6',
    'user6@example.com',
    'User Six',
    crypt('password6', gen_salt('bf')),
    '404 Cedar St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user7',
    'user7@example.com',
    'User Seven',
    crypt('password7', gen_salt('bf')),
    '505 Walnut St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user8',
    'user8@example.com',
    'User Eight',
    crypt('password8', gen_salt('bf')),
    '606 Chestnut St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user9',
    'user9@example.com',
    'User Nine',
    crypt('password9', gen_salt('bf')),
    '707 Ash St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  ),
  (
    'user10',
    'user10@example.com',
    'User Ten',
    crypt('password10', gen_salt('bf')),
    '808 Poplar St',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
  );

COMMIT;