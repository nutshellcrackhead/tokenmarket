package services

var bonusQueries map[string]string = map[string]string{
	/**
	 * Put bonuses from deposits
	 *
	 * @param {int} time - time now
	 */
	"PutDepositsBonuses": `SELECT FROM deposits_bonuses($1);`,

	/**
	 * Put bonuses from referral
	 *
	 * @param {int} time - time now
	 * @param {int} last_bonus - time when referral bonus was put last time
	 */
	"PutReferralBonuses": `SELECT FROM referral_bonuses($1, $2);`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return bonusQueries
}
