package services

var walletQueries map[string]string = map[string]string{
	/**
	 * Returns user's wallet status
	 *
	 * @param {int} id - user id
	 * @returns {{ float currency, int amount }} - salt and password
	 */
	"GetWalletStatus": `SELECT currency, amount FROM token_wallets
		WHERE username = $1;`,

	/**
	 * Returns user's operations list
	 *
	 * @param {int} id - user id
	 * @param {int} limit - results on one page
	 * @param {int} offset - start from result number
	 * @returns {{ float currency, int amount }} - salt and password
	 */
	"GetUserOperations": `(SELECT token_operations.id, date, status, parse_topup_type(token_topup.method) AS "type",
			token_topup.amount,
			'income' AS action,
			token_topup.currency
		FROM token_operations
		LEFT OUTER JOIN token_topup
		ON token_operations.id = token_topup.operation
		WHERE (type = 'token_topup' AND token_operations.user = $1))
		UNION ALL
			(SELECT token_operations.id, date, status, token_operations.type,
				token_topup.amount,
				'expense' AS action,
				token_topup.currency
			FROM token_operations
			LEFT OUTER JOIN token_topup
			ON token_operations.id = token_topup.operation
			WHERE (type = 'token_account_activation' AND token_operations.user = $1))
		UNION ALL
			(SELECT token_operations.id, date, status, token_operations.type,
				token_deposits_create.amount,
				'expense' AS action,
				token_deposits_create.currency
			FROM token_operations
			LEFT OUTER JOIN token_deposits_create
			ON token_operations.id = token_deposits_create.operation
			WHERE (type = 'token_deposits_create' AND token_operations.user = $1))
		UNION ALL
			(SELECT token_operations.id, date, status, token_operations.type,
				token_payout_operation.amount,
				'expense' AS action,
				token_payout_operation.currency
			FROM token_operations
			LEFT OUTER JOIN token_payout_operation
			ON token_operations.id = token_payout_operation.operation
			WHERE (type = 'token_payout_operation' AND token_operations.user = $1))
		UNION ALL
			(SELECT token_operations.id, date, status, token_operations.type,
				token_deposits_bonus.amount,
				'income' AS action,
				token_deposits_bonus.currency
			FROM token_operations
			LEFT OUTER JOIN token_deposits_bonus
			ON token_operations.id = token_deposits_bonus.operation
			WHERE (type = 'token_deposits_bonus' AND token_operations.user = $1))
		UNION ALL
			(SELECT token_operations.id, date, status, token_operations.type,
				token_referral_bonus.amount,
				'income' AS action,
				token_referral_bonus.currency
			FROM token_operations
			LEFT OUTER JOIN token_referral_bonus
			ON token_operations.id = token_referral_bonus.operation
			WHERE (type = 'token_referral_bonus' AND token_operations.user = $1))
		ORDER BY date DESC
		LIMIT $2 OFFSET $3;`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return walletQueries
}
