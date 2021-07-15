package internal

//Config struct to set flags and connection info
type Config struct {
	LogFormat     string
	LogLevel      string
	ListenAddress string
	ListenPath    string
	SentryURI     string
	SentryAPIKey  string
	SentryOrg     string
	Interval      string
}
