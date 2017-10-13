package services

var operationsQueries map[string]string = map[string]string{
	/**
	 * Gets referrers list by referrer id
	 *
	 * @param {int} id - referrer id
	 * @param {int} limit - results per page
	 * @param {int} offset - offset of results
	 * @returns {{}}
	 */
	"CreatePayment": `WITH payment_operation AS
		(INSERT INTO "public"."token_operations"("user", "type", "status", "date")
			VALUES($1, $7, DEFAULT, $3) RETURNING "id")
		INSERT INTO "public"."token_topup" ("operation", "method", "amount", "currency", "token")
		SELECT payment_operation.id, $2, $4, $5, $6
		FROM payment_operation
		RETURNING
			"id",
			"token",
			"amount",
			"currency",
			"method";`,

	"CreatePayout": `SELECT FROM create_payout($1, $2, $3, $4, $5, $6, $7, $8);`,

	"CreateAccountActivation": `INSERT INTO "public"."token_account_activation"("payment")
		VALUES($1);`,

	"ValidatePayment": `SELECT * FROM validate_topup_status($1, $2, $3, $4, $5);`,

	"UpdatePayment": `UPDATE token_operations
		SET status = $2
		WHERE id = $1;`,

	"UnlockAccount": `SELECT activate_user($1);`,

	"TopUpWallet": `SELECT topup_wallet($1)`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return operationsQueries
}
