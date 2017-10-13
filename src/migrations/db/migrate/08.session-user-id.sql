ALTER TABLE "public"."token_sessions" ADD FOREIGN KEY ("username") REFERENCES "public"."token_users"("id");
