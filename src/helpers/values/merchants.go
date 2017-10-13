package values

var merchants map[string]string = map[string]string{
	"perfectmoney": "U12258828",
	"advcash":      "tokenmarket1@gmail.com",
}

func GetMerchant(merchant string) string {
	return merchants[merchant]
}
