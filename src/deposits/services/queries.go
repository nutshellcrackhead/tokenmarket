package services

var depositsQueries map[string]string = map[string]string{
	/**
	 * Returns user's wallet status
	 *
	 * @param {int} id - user id
	 * @param {int} page - page
	 * @returns {{ int date, int valid, string currency, float amount,
	 *	string status, int paidout, int paydate }} - collection of deposits
	 */
	"GetDepositsList": `SELECT "date", "valid", "currency", "amount",
		"status", "paidout", "paydate" FROM token_deposits
		WHERE user_id = $1
		LIMIT $2 OFFSET $3;`,

	/**
	 * Creates deposits and returns new deposit instance or raise an error
	 *
	 * @param {int} id - user id
	 * @param {float} deposit_amount - deposit amount
	 * @param {string} deposit_currency - deposit currency
	 * @param {int} operation_date - time now in unix nano
	 * @param {int} valid_till - unix nano time till when deposit is active
	 * @returns {{ int date, int valid, string currency, float amount,
	 *	string status, int paidout, int paydate }} - deposit
	 */
	"CreateDeposit": `SELECT "date", "valid", "currency", "amount", "status",
	 	"paidout", "paydate" FROM create_deposit($1, $2, $3, $4, $5);`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return depositsQueries
}
