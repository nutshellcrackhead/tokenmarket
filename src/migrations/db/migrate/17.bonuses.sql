CREATE OR REPLACE FUNCTION calculate_deposit_bonus(deposit token_deposits, time_now bigint) RETURNS float AS $$
DECLARE
  till_date bigint;
  since_date bigint;
  bonus_per_day float;
  amount_to_pay float;
  day_in_nanosecs bigint;
  days_to_pay float;
BEGIN
  day_in_nanosecs = 86400000065018;
  bonus_per_day = deposit.amount * 0.15 * 12 / 365;
  till_date = time_now;
  amount_to_pay = 0;
  since_date = deposit.paydate;

  IF deposit.valid < time_now THEN
    UPDATE token_deposits SET status = 'expired' WHERE id = deposit.id;
    till_date = deposit.valid;
    amount_to_pay = amount_to_pay + deposit.amount;
  END IF;

  days_to_pay = ((till_date - since_date) / day_in_nanosecs);
  amount_to_pay = amount_to_pay + days_to_pay * bonus_per_day;
  UPDATE token_deposits SET paidout = token_deposits.paidout + amount_to_pay, paydate = time_now
  WHERE id = deposit.id;

  RETURN amount_to_pay;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION deposits_bonuses(time_now bigint) RETURNS VOID AS $$
DECLARE
  deposit token_deposits;
  bonus float;
BEGIN
  FOR deposit IN (SELECT * FROM token_deposits WHERE token_deposits.valid > token_deposits.paydate AND status != 'expired') LOOP
    SELECT * INTO bonus FROM calculate_deposit_bonus(deposit, time_now);
    PERFORM create_deposit_bonus_operation(deposit.user_id, time_now, bonus, deposit.currency);
    UPDATE token_wallets SET amount = token_wallets.amount + bonus WHERE currency = deposit.currency AND username = deposit.user_id;
  END LOOP;
END;
$$ LANGUAGE plpgsql;
