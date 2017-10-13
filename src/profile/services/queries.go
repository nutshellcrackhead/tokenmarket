package services

var profileQueries map[string]string = map[string]string{
	/**
	 * Gets full profile by id
	 *
	 * @param {integer} id - user id
	 * @returns {Profile}
	 */
	"GetProfile": `SELECT id, username, avatar, skype, email, phone, muted,
			token_locations.city_id "token_locations.id",
			token_countries.title "token_locations.country",
			token_regions.title "token_locations.region",
			token_locations.title "token_locations.city"
		FROM token_users
		LEFT OUTER JOIN token_locations ON (token_users.location = token_locations.city_id)
		LEFT OUTER JOIN token_countries ON (token_locations.country_id = token_countries.country_id)
		LEFT OUTER JOIN token_regions ON (token_locations.region_id = token_regions.region_id)
		WHERE id = $1;`,

	/**
	 * Updates profile data by id
	 * Doesn't update data if provided parameter is null
	 *
	 * @param {integer} id - user id
	 * @param {string} name - name
	 * @param {integer ref token_locations city_id} location - location id
	 * @param {string} skype - skype id
	 * @param {string} avatar - avatar url
	 * @returns {Profile} - updated profile
	 */
	"UpdateProfileData": `WITH update_user
		AS (UPDATE token_users SET
			location = COALESCE($2, location),
			skype = COALESCE($3, skype),
			avatar = COALESCE($4, avatar)
			WHERE "id"=$1
			RETURNING *
		)
		SELECT update_user.id,
			update_user.username,
			update_user.skype,
			update_user.phone,
			update_user.email,
			update_user.avatar,
			update_user.muted,
			token_locations.city_id "token_locations.id",
			token_countries.title "token_locations.country",
			token_regions.title "token_locations.region",
			token_locations.title "token_locations.city"
		FROM update_user
		LEFT OUTER JOIN token_locations ON (update_user.location = token_locations.city_id)
		LEFT OUTER JOIN token_countries ON (token_locations.country_id = token_countries.country_id)
		LEFT OUTER JOIN token_regions ON (token_locations.region_id = token_regions.region_id);`,

	/**
	 * Gets cities list by name fragment
	 *
	 * @param {string} nameFragment - city name fragment
	 * @returns {City}
	 */
	"GetCitiesList": `SELECT
			city_id AS id,
			token_locations.title AS city,
			token_regions.title AS region,
			token_countries.title AS country
		FROM token_locations
		LEFT OUTER JOIN token_countries
		ON (token_locations.country_id = token_countries.country_id)
		LEFT OUTER JOIN token_regions
		ON (token_locations.region_id = token_regions.region_id)
		WHERE LOWER(token_locations.title) LIKE LOWER($1 || '%')
		ORDER BY LOWER(token_locations.title) ASC
		LIMIT 10;`,

	/**
	 * Gets city by id
	 *
	 * @param {string} id - city id
	 * @returns {City}
	 */
	"GetCity": `SELECT city_id AS id,
			token_locations.title AS city,
			token_regions.title AS region,
			token_countries.title AS country
		FROM token_locations
		LEFT OUTER JOIN token_countries
		ON (token_locations.country_id = token_countries.country_id)
		LEFT OUTER JOIN token_regions
		ON (token_locations.region_id = token_regions.region_id)
		WHERE city_id = $1;`,
}

// Returns map with Postgres queries
func GetQueries() map[string]string {
	return profileQueries
}
