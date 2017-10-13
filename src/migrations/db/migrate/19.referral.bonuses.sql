CREATE OR REPLACE FUNCTION create_referral_bonus_operation(user_id int, time_now bigint, amount float, currency token_wallet_currencies) RETURNS VOID AS $$
BEGIN
  WITH payment_operation AS
  (INSERT INTO "public"."token_operations"("user", "type", "status", "date")
    VALUES(user_id, 'token_referral_bonus', 'success', time_now) RETURNING "id")
  INSERT INTO "public"."token_referral_bonus" ("operation", "amount", "currency")
    SELECT payment_operation.id, amount, currency
    FROM payment_operation;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION referral_bonuses(time_now bigint, time_last_bonus bigint) RETURNS VOID AS $$
DECLARE
  deposit token_deposits;
  leg_name token_referral_leg;
  working_user_id int;
  tree_id int;
BEGIN
  FOR deposit IN (SELECT * FROM token_deposits WHERE date >= time_last_bonus AND date < time_now) LOOP
    UPDATE token_revenue SET revenue = revenue + deposit.amount WHERE user_id = deposit.user_id AND currency = deposit.currency;

    working_user_id = deposit.user_id;

    SELECT tree INTO tree_id FROM token_referral WHERE user_id = working_user_id;

    WHILE working_user_id <> tree_id LOOP

      SELECT parent, leg INTO working_user_id, leg_name FROM token_referral WHERE user_id = working_user_id;

      IF leg_name = 'left' THEN
        UPDATE token_revenue SET "left" = "left" + deposit.amount, date = time_now WHERE user_id = working_user_id AND currency = deposit.currency;
      END IF;

      IF leg_name = 'right' THEN
        UPDATE token_revenue SET "right" = "right" + deposit.amount, date = time_now WHERE user_id = working_user_id AND currency = deposit.currency;
      END IF;
    END LOOP;
  END LOOP;

  PERFORM checkout_revenue(time_now);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION checkedout_bonus(checkedout_user_id int, checkedout_amount float, checkedout_currency token_wallet_currencies, time_now bigint) RETURNS VOID AS $$
DECLARE
  checkedout_actual float;
  bonus_amount float;
  checkedout_sum float;
BEGIN
  bonus_amount = 0;

  SELECT checkedout INTO checkedout_actual FROM token_revenue WHERE "token_revenue"."user_id" = checkedout_user_id AND currency = checkedout_currency;

  checkedout_sum = checkedout_actual + checkedout_amount;

  UPDATE token_revenue SET checkedout = checkedout + checkedout_amount WHERE "user_id" = checkedout_user_id AND currency = checkedout_currency;

  IF checkedout_currency <> 'USD' THEN
    RETURN;
  END IF;

  IF (checkedout_actual < 1000 AND checkedout_sum >= 1000) THEN
    bonus_amount = bonus_amount + 50;
  END IF;

  IF (checkedout_actual < 5000 AND checkedout_sum >= 5000) THEN
    bonus_amount = bonus_amount + 250;
  END IF;

  IF (checkedout_actual < 15000 AND checkedout_sum >= 15000) THEN
    bonus_amount = bonus_amount + 750;
  END IF;

  IF (checkedout_actual < 50000 AND checkedout_sum >= 50000) THEN
    bonus_amount = bonus_amount + 2500;
  END IF;

  IF (checkedout_actual < 150000 AND checkedout_sum >= 150000) THEN
    bonus_amount = bonus_amount + 7500;
  END IF;

  IF (checkedout_actual < 500000 AND checkedout_sum >= 500000) THEN
    bonus_amount = bonus_amount + 25000;
  END IF;

  IF (checkedout_actual < 1000000 AND checkedout_sum >= 1000000) THEN
    bonus_amount = bonus_amount + 40000;
  END IF;

  IF (bonus_amount = 0 ) THEN
    RETURN;
  END IF;

  PERFORM create_referral_bonus_operation(checkedout_user_id, time_now, bonus_amount, checkedout_currency);

  UPDATE token_wallets SET amount = amount + bonus_amount
  WHERE "username" = checkedout_user_id AND currency = checkedout_currency;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION checkout_revenue(time_now bigint) RETURNS VOID AS $$
DECLARE
  user_revenue token_revenue;
  checkedout_amount float;
  amount_left float;
  amount_right float;
BEGIN
  FOR user_revenue IN (SELECT * FROM token_revenue) LOOP
    checkedout_amount = 0;
    amount_left = 0;
    amount_right = 0;

    IF (user_revenue.left = 0 OR user_revenue.right = 0) THEN
      CONTINUE;
    END IF;

    IF (user_revenue.left > user_revenue.right) THEN
      checkedout_amount = user_revenue.right;
      amount_left = user_revenue.left - user_revenue.right;
    ELSE
      checkedout_amount = user_revenue.left;
      amount_right = user_revenue.right - user_revenue.left;
    END IF;

    UPDATE token_revenue SET "left" = amount_left, "right" = amount_right
    WHERE id = user_revenue.id;

    PERFORM create_referral_bonus_operation(user_revenue.user_id, time_now, (checkedout_amount * 0.06), user_revenue.currency);

    PERFORM checkedout_bonus(user_revenue.user_id, checkedout_amount, user_revenue.currency, time_now);

    UPDATE token_wallets SET amount = amount + (checkedout_amount * 0.06)
    WHERE "username" = user_revenue.user_id AND currency = user_revenue.currency;
  END LOOP;
END;
$$ LANGUAGE plpgsql;
