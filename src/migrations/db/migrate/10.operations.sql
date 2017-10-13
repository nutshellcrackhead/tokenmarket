CREATE TYPE token_wallet_currencies AS ENUM('USD');

ALTER TABLE "public"."token_users"
  ALTER COLUMN "salt" SET NOT NULL,
  ALTER COLUMN "password" SET NOT NULL;

CREATE TYPE token_payment_method AS ENUM('perfectmoney', 'advcash', 'payeer', 'blockio');
CREATE TYPE token_operation_type AS ENUM('token_phone_confirmation', 'token_topup',
  'token_account_activation',
  'token_pm_topup',
  'token_payeer_topup',
  'token_advcash_topup',
  'token_blockio_topup',
  'token_deposits_create',
  'token_deposits_bonus',
  'token_payout_operation',
  'token_referral_bonus');

CREATE OR REPLACE FUNCTION parse_topup_type(payment_method token_payment_method) RETURNS token_operation_type AS $$
BEGIN
  IF payment_method = 'perfectmoney' THEN
    RETURN 'token_pm_topup';
  END IF;

  IF payment_method = 'advcash' THEN
    RETURN 'token_advcash_topup';
  END IF;

  IF payment_method = 'blockio' THEN
    RETURN 'token_blockio_topup';
  END IF;

  IF payment_method = 'payeer' THEN
    RETURN 'token_payeer_topup';
  END IF;
END;
$$ LANGUAGE plpgsql;


CREATE TYPE token_operation_status AS ENUM('on_hold', 'in_process', 'timeout', 'failed',
  'cancelled', 'success');

CREATE TYPE token_operation_transaction_values AS ENUM('income', 'expense');

CREATE TABLE "public"."token_operations" (
  "id" serial,
  "date" bigint NOT NULL,
  "status" token_operation_status NOT NULL DEFAULT 'in_process',
  "type" token_operation_type NOT NULL,
  "user" int,
  PRIMARY KEY ("id"),
  UNIQUE ("id"),
  FOREIGN KEY ("user") REFERENCES "public"."token_users"("id")
);

CREATE TABLE "public"."token_topup" (
  "id" serial,
  "operation" int,
  "method" token_payment_method NOT NULL,
  "amount" float NOT NULL,
  "currency" token_wallet_currencies NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id")
);

CREATE TABLE "public"."token_account_activation" (
  "id" serial,
  "payment" int,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("payment") REFERENCES "public"."token_topup"("id")
);

CREATE TABLE "public"."token_operations_confirmation" (
  "id" serial,
  "operation" int NOT NULL,
  "token" bytea NOT NULL,
  "valid_till" bigint NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id")
);

CREATE TABLE "public"."token_phone_confirmation" (
  "id" serial,
  "operation" int NOT NULL,
  "confirmation" int NOT NULL,
  "new_phone" text NOT NULL,
  "new_token" bytea NOT NULL,
  "confirmed" boolean NOT NULL DEFAULT 'false',
  PRIMARY KEY ("id"),
  FOREIGN KEY ("operation") REFERENCES "public"."token_operations"("id"),
  FOREIGN KEY ("confirmation") REFERENCES "public"."token_operations_confirmation"("id")
);

