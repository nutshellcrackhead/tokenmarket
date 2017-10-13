package services

var referralQueries map[string]string = map[string]string{
	/**
	 * Gets referrers list by referrer id
	 *
	 * @param {int} id - referrer id
	 * @param {int} limit - results per page
	 * @param {int} offset - offset of results
	 * @returns {{}}
	 */
	"GetReferrers": `SELECT token_users.username,
			token_users.avatar,
			level,
			leg,
			token_users.registered,
			revenue_user.revenue AS "revenue",
			referrer_user.username AS "referrer",
			parent_user.username AS "parent"
		FROM get_referrers($1, NULL, NULL)
			LEFT OUTER JOIN token_users
			ON (token_users.id = user_id)
			LEFT OUTER JOIN token_users AS referrer_user
			ON (referrer_user.id = referrer)
			LEFT OUTER JOIN token_users AS parent_user
			ON (parent_user.id = parent)
			LEFT OUTER JOIN token_revenue AS revenue_user
			ON (revenue_user.user_id = token_users.id)
		ORDER BY "level" ASC LIMIT $2 OFFSET $3;`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return referralQueries
}
