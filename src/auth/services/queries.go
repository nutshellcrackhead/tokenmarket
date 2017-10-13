package services

var authQueries map[string]string = map[string]string{
	/**
	 * Returns user's salt and password by username
	 *
	 * @param {string} username - username
	 * @returns {{ string salt, string password }} - salt and password
	 */
	"GetUserCredentials": `SELECT id, salt, password
		FROM token_users
		WHERE username = LOWER($1);`,

	/**
	 * Registers user with provided data
	 *
	 * @param {string} username - username
	 * @param {string} phone - phone
	 * @param {string} email - email
	 * @param {[]byte} salt - salt
	 * @param {[]byte} password - password
	 * @returns {int64} - new user's id
	 */
	"RegisterUser": `INSERT INTO
		"public"."token_users"("username", "phone", "email", "salt", "password", "registered")
		VALUES(LOWER($1), $2, $3, $4, $5, $6)
		RETURNING id;`,

	/**
	 * Save user's session in the database
	 *
	 * @param {int64} user - user's id
	 * @param {string} token - session's token
	 * @param {int64} valid_till - valid till timestamp in unix nano
	 * @returns {int64, int64, string, time.Timestamp} - new session object
	 */
	"SaveSessionToken": `INSERT INTO
		"public"."token_sessions"("username", "token", "valid_till")
		VALUES($1, $2, $3);`,

	/**
	 * Invalidates all the opened sessions for specific user
	 *
	 * @param {int64} user - username
	 * @param {int64} valid_till - valid till timestamp in unix nano
	 */
	"InvalidateOtherSessions": `UPDATE token_sessions SET valid_till = $2
		WHERE (username = $1 AND valid_till > $2);`,

	/**
	 * Prolongs user's session defined by token
	 *
	 * @param {string} token - token
	 * @param {int64} valid_till - valid till timestamp in unix nano
	 */
	"ProlongSession": `UPDATE token_sessions SET valid_till = $2
		WHERE (token = $1);`,

	/**
	 * Creates wallets for each value of token_wallet_currencies. If wallet for specific
	 * currency exists - creation of the wallet is omitted
	 *
	 * @param {int} id - user id
	 */
	"CreateWallets": `SELECT create_wallets($1);`,

	/**
	 * Checks if referrer exists
	 * @param {string} username - name of the user that is checked
	 * @returns {bool}
	 */
	"ValidateReferrer": `SELECT id FROM token_users
			WHERE username = LOWER($1);`,

	/**
	 * Registers referrers
	 * @param {int} user_id - user id to save
	 * @param {int} referrer_id - referrer id
	 * @param {string} leg - left or right
	 */
	"RegisterReferrer": `SELECT save_referrer($1, $2, $3)`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return authQueries
}
