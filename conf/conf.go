package conf

const (
	WEB_ADDRESS = "localhost"
	USER_PORT   = "8080"
	//Dont store the signing key in script, use environment variables to get the signing key for jwt generation.
	Signing_Key = "SUPERSECRETKEY"
	TIME_LAYOUT = "2006-01-02T15:05:05"
)

func GetServerAddress() string {
	return WEB_ADDRESS + ":" + USER_PORT
}
