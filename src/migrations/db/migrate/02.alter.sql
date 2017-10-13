-- ste autoincrement id field
ALTER TABLE token_users ALTER COLUMN id SET DEFAULT nextval('token_users_id_seq');