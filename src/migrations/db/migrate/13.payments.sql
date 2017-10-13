ALTER TABLE "public"."token_topup" ADD COLUMN "token" text NOT NULL;

CREATE OR REPLACE FUNCTION validate_topup_status(payment_id integer, token text, amount float, currency token_wallet_currencies, method_name token_payment_method) RETURNS int AS $$
DECLARE
  id_topup int;
  operation_topup int;
  token_topup text;
  amount_topup float;
  currency_topup token_wallet_currencies;
  method_name_topup token_payment_method;
  payment_status token_operation_status;
BEGIN
  SELECT token_topup.id, token_topup.token, token_topup.amount, token_topup.currency, token_topup.method, token_topup.operation
  INTO id_topup, token_topup, amount_topup, currency_topup, method_name_topup, operation_topup
  FROM token_topup
  WHERE token_topup.id = payment_id;

  SELECT "token_operations"."status"
  INTO payment_status
  FROM token_operations
  WHERE id = operation_topup;

  IF (payment_status = 'success' OR payment_status = 'failed') THEN
    RETURN 0;
  END IF;

  IF (payment_id = id_topup AND token = token_topup AND amount = amount_topup AND currency = currency_topup AND method_name = method_name_topup) THEN
    RETURN operation_topup;
  END IF;

  RETURN 0;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION activate_user(operation_id int) RETURNS VOID AS $$
DECLARE
  user_id int;
  status text;
  operation_type token_operation_type;
BEGIN
  SELECT token_operations.user, token_operations.status, token_operations.type INTO user_id, status, operation_type FROM token_operations
  WHERE token_operations.id = operation_id;

  IF (operation_type = 'token_account_activation' AND status = 'success') THEN
    UPDATE token_users SET muted = false WHERE id = user_id;
  END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION topup_wallet(operation_id integer) RETURNS VOID AS $$
DECLARE
  amount_topup float;
  currency_topup token_wallet_currencies;
  username_topup int;
  operation_type token_operation_type;
BEGIN
  SELECT amount, currency, "user", token_operations.type
  INTO amount_topup, currency_topup, username_topup, operation_type
  FROM token_topup
    LEFT OUTER JOIN token_operations ON token_operations.id = operation_id
  WHERE operation = operation_id;

  IF (amount_topup > 0 AND operation_type = 'token_topup') THEN
    UPDATE token_wallets SET amount = amount + amount_topup WHERE (token_wallets.username = username_topup AND token_wallets.currency = currency_topup);
  END IF;
END;
$$ LANGUAGE plpgsql;
