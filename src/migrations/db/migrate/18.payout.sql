CREATE TABLE "public"."token_payout_operation" (
  "id" serial,
  "amount" text,
  "currency" text,
  "operation" int,
  "method" token_payment_method,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id")
);

ALTER TABLE "public"."token_payout_operation"
  ALTER COLUMN "amount" SET DATA TYPE float USING amount::double precision,
  ADD COLUMN "token" text NOT NULL,
  ALTER COLUMN "amount" SET NOT NULL,
  ALTER COLUMN "currency" SET NOT NULL,
  ALTER COLUMN "method" SET NOT NULL,
  ADD COLUMN "account" text NOT NULL;

CREATE OR REPLACE FUNCTION create_payout(user_id int, payout_method token_payment_method, operation_date bigint, payout_amount float, payout_currency token_wallet_currencies, payout_token text, account_credential text, operation_type token_operation_type) RETURNS VOID AS $$
DECLARE
  wallet_balance float;
BEGIN
  SELECT "token_wallets"."amount" INTO wallet_balance FROM token_wallets WHERE username = user_id AND currency = payout_currency;

  IF (payout_amount <= 0 OR payout_amount > wallet_balance) THEN
    WITH payout_operation AS
    (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
      VALUES(user_id, operation_type, 'failed', operation_date) RETURNING "id")
    INSERT INTO "public"."token_payout_operation" ("operation", "method", "amount", "currency", "token", "account")
      SELECT payout_operation.id, payout_method, payout_amount, payout_currency, payout_token, account_credential
      FROM payout_operation;

    RAISE EXCEPTION 'Invalid amount';
    RETURN;
  END IF;


  WITH payout_operation AS
  (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
    VALUES(user_id, operation_type, DEFAULT, operation_date) RETURNING "id")
  INSERT INTO "public"."token_payout_operation" ("operation", "method", "amount", "currency", "token", "account")
    SELECT payout_operation.id, payout_method, payout_amount, payout_currency, payout_token, account_credential
    FROM payout_operation;

  UPDATE token_wallets SET amount = token_wallets.amount - payout_amount WHERE username = user_id AND currency = payout_currency;
END;
$$ LANGUAGE plpgsql;

ALTER TABLE "public"."token_payout_operation" ALTER COLUMN "currency" SET DATA TYPE token_wallet_currencies USING currency::token_wallet_currencies;
