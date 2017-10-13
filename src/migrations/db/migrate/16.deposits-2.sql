CREATE TABLE "public"."token_deposits_create" (
  "id" serial,
  "amount" float DEFAULT 0,
  "currency" token_wallet_currencies DEFAULT 'USD',
  "operation" int NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id"),
  UNIQUE ("operation")
);

CREATE TABLE "public"."token_deposits_bonus" (
  "id" serial,
  "operation" int NOT NULL,
  "amount" float NOT NULL,
  "currency" token_wallet_currencies DEFAULT 'USD',
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id")
);

CREATE TABLE "public"."token_referral_bonus" (
  "id" serial,
  "operation" int NOT NULL,
  "amount" float NOT NULL,
  "currency" token_wallet_currencies DEFAULT 'USD',
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id")
);

CREATE OR REPLACE FUNCTION create_deposit_bonus_operation(user_id int, time_now bigint, payout float, currency token_wallet_currencies) RETURNS VOID AS $$
BEGIN
  WITH payment_operation AS
  (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
    VALUES(user_id, 'token_deposits_bonus', 'success', time_now) RETURNING "id")
  INSERT INTO "public"."token_deposits_bonus" ("operation", "amount", "currency")
    SELECT payment_operation.id, payout, currency
    FROM payment_operation;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_deposit(deposit_user_id int, deposit_amount float, deposit_currency token_wallet_currencies, operation_date bigint, valid_till bigint) RETURNS SETOF token_deposits AS $$
DECLARE
  wallet_balance float;
  referrer_id int;
BEGIN
  SELECT "token_wallets"."amount" INTO wallet_balance FROM token_wallets WHERE "token_wallets"."username" = deposit_user_id AND currency = deposit_currency;

  IF (deposit_amount <= 0 OR deposit_amount > wallet_balance) THEN
    WITH deposit_operation AS
    (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
      VALUES(deposit_user_id, 'token_deposits_create', 'failed', operation_date) RETURNING "id")
    INSERT INTO "public"."token_deposits_create" ("operation", "amount", "currency")
      SELECT deposit_operation.id, deposit_amount, deposit_currency
      FROM deposit_operation;

    RAISE EXCEPTION 'Invalid amount';
    RETURN;
  END IF;

  UPDATE token_wallets SET amount = token_wallets.amount - deposit_amount WHERE username = deposit_user_id AND currency = deposit_currency;

  WITH deposit_operation AS
  (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
    VALUES(deposit_user_id, 'token_deposits_create', 'success', operation_date) RETURNING "id")
  INSERT INTO "public"."token_deposits_create" ("operation", "amount", "currency")
    SELECT deposit_operation.id, deposit_amount, deposit_currency
    FROM deposit_operation;

  SELECT referrer INTO referrer_id FROM token_referral WHERE "token_referral"."user_id" = deposit_user_id;

  PERFORM create_deposit_bonus_operation(referrer_id, operation_date, (deposit_amount * 0.1), deposit_currency);

  UPDATE token_wallets SET amount = token_wallets.amount + (deposit_amount * 0.1)
  WHERE token_wallets.username = referrer_id AND token_wallets.currency = deposit_currency;

  RETURN QUERY INSERT INTO "token_deposits"("amount", "currency", "date", "valid", "status", "user_id", "paidout", "paydate")
  VALUES(deposit_amount, deposit_currency, operation_date, valid_till, DEFAULT, deposit_user_id, DEFAULT, operation_date) RETURNING *;
END;
$$ LANGUAGE plpgsql;

ALTER TABLE "public"."token_deposits" ALTER COLUMN "paydate" SET DATA TYPE bigint;
