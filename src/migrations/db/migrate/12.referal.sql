CREATE TYPE token_referral_leg AS ENUM('left', 'right');

CREATE TABLE "public"."token_referral" (
  "id" serial,
  "user" int NOT NULL,
  "referrer" int,
  "parent" int,
  "level" int DEFAULT 0,
  "tree" int NOT NULL,
  "leg" token_referral_leg,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user") REFERENCES "public"."token_users"("id"),
  FOREIGN KEY ("referrer") REFERENCES "public"."token_users"("id"),
  FOREIGN KEY ("parent") REFERENCES "public"."token_users"("id"),
  FOREIGN KEY ("tree") REFERENCES "public"."token_users"("id")
);

ALTER TABLE "public"."token_users" ADD COLUMN "registered" bigint;

ALTER TABLE "public"."token_referral"
  ALTER COLUMN "level" SET NOT NULL;

ALTER TABLE "public"."token_referral" RENAME COLUMN "user" TO "user_id";

CREATE TABLE "public"."token_revenue" (
  "id" serial,
  "user_id" int NOT NULL,
  "currency" token_wallet_currencies DEFAULT 'USD',
  "left" float DEFAULT 0,
  "right" float DEFAULT 0,
  "date" bigint DEFAULT 0,
  "revenue" float DEFAULT 0,
  "checkedout" float DEFAULT 0,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("user_id") REFERENCES "public"."token_users"("id")
);


CREATE OR REPLACE FUNCTION get_referrers(referrer_id int, parent_leg token_referral_leg, parent_level int) RETURNS setof token_referral AS $$
DECLARE
  referral token_referral;
  parent_leg_temp token_referral_leg;
  temp_parent_level int;
BEGIN
  temp_parent_level = COALESCE(parent_level, 1);

  FOR referral IN (SELECT id, user_id, referrer, parent, temp_parent_level AS "level", tree, COALESCE(parent_leg, leg) AS "leg"
                   FROM token_referral WHERE token_referral.parent = referrer_id) LOOP

    IF parent_leg IS NULL THEN
      parent_leg_temp = referral.leg;
    ELSE
      parent_leg_temp = parent_leg;
    END IF;

    RETURN NEXT referral;

    RETURN QUERY (SELECT * FROM get_referrers(referral.user_id, parent_leg_temp, temp_parent_level + 1));
  END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION save_referrer(save_user_id int, referrer_id int, leg_name token_referral_leg) RETURNS VOID AS $$
DECLARE
  level_num int;
  tree_id int;
  parent_id int;
  temp_parent_id int;
BEGIN
  SELECT token_referral.level, token_referral.tree
  INTO level_num, tree_id
  FROM token_referral
  WHERE user_id = referrer_id;

  temp_parent_id = referrer_id;

  WHILE temp_parent_id IS NOT NULL LOOP
    level_num = level_num + 1;
    parent_id = temp_parent_id;

    SELECT user_id INTO temp_parent_id FROM token_referral WHERE
      (token_referral.tree = tree_id AND
       token_referral.level = level_num AND
       token_referral.parent = temp_parent_id AND
       token_referral.leg = leg_name);
  END LOOP;

  INSERT INTO "public"."token_referral"("user_id", "referrer", "parent", "level", "tree", "leg")
  VALUES (save_user_id, referrer_id, parent_id,
          COALESCE(level_num, 0),
          COALESCE(tree_id, save_user_id),
          leg_name);

  INSERT INTO token_revenue
  VALUES (DEFAULT, save_user_id, 'USD', DEFAULT, DEFAULT, DEFAULT, DEFAULT);
END;
$$ LANGUAGE plpgsql;
