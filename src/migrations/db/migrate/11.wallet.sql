CREATE TABLE "public"."token_wallets" (
  "id" serial,
  "username" integer NOT NULL,
  "amount" float NOT NULL DEFAULT 0,
  "currency" token_wallet_currencies NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("username") REFERENCES "public"."token_users"("id")
);

ALTER TABLE "public"."token_users"
  ALTER COLUMN "muted" SET DATA TYPE boolean
  USING isfinite(muted),
  ALTER COLUMN "muted" SET DEFAULT true;

ALTER TABLE "public"."token_users" ALTER COLUMN "muted" SET NOT NULL;

CREATE OR REPLACE FUNCTION create_wallets(user_id int) RETURNS void AS $$
DECLARE
  currency_name   token_wallet_currencies;
  existing_wallet int;
BEGIN
  FOREACH currency_name IN ARRAY ARRAY(SELECT enum_range(NULL::token_wallet_currencies))
  LOOP
    existing_wallet = (SELECT id FROM "public"."token_wallets" WHERE "username" = user_id AND "currency" = currency_name);
    IF existing_wallet IS NULL THEN
      INSERT INTO "public"."token_wallets"("username", "currency") VALUES (user_id, currency_name);
    END IF;
  END LOOP;
END;
$$ LANGUAGE plpgsql;
