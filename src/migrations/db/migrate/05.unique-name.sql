ALTER TABLE "public"."token_users"
  ADD UNIQUE ("name");
ALTER TABLE "public"."token_users" RENAME COLUMN "name" TO "username";