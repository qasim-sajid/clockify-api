package conf

const (
	WEB_ADDRESS = "localhost"
	USER_PORT   = "8080"
)

func GetServerAddress() string {
	return WEB_ADDRESS + ":" + USER_PORT
}
