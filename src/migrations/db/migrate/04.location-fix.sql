ALTER TABLE "public"."token_users"
  ALTER COLUMN "location" SET DATA TYPE SMALLINT USING location::SMALLINT,
  ADD FOREIGN KEY ("location") REFERENCES "public"."token_locations"("city_id"),
  ADD UNIQUE ("phone");
