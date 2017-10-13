ALTER TABLE "public"."token_users"
  ALTER COLUMN salt DROP DEFAULT,
  ALTER COLUMN salt SET DATA TYPE bytea USING decode('', 'escape'),
  ALTER COLUMN password DROP DEFAULT,
  ALTER COLUMN password SET DATA TYPE bytea USING decode('', 'escape');
