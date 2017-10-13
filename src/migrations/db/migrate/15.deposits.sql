CREATE TYPE token_deposit_status AS ENUM('valid', 'expired');

CREATE TABLE "public"."token_deposits" (
  "id" serial,
  "amount" float NOT NULL,
  "currency" token_wallet_currencies NOT NULL DEFAULT 'USD',
  "date" bigint NOT NULL,
  "valid" bigint NOT NULL,
  "status" token_deposit_status NOT NULL DEFAULT 'valid',
  "user_id" int,
  "paidout" float NOT NULL DEFAULT 0,
  "paydate" int NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."token_users"("id")
);