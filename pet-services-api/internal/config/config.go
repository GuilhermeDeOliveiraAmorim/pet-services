package config

import "os"

type DB struct {
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
	DB_NAME     string
}

type FRONT_END_URL struct {
	FRONT_END_URL_PROD string
	FRONT_END_URL_DEV  string
}

type SECRETS struct {
	JWT_SECRET         string
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
}

type GOOGLE struct {
	IMAGE_BUCKET_NAME string
	URL_BUCKET_NAME   string
}

type AVAILABLE_LANGUAGES struct {
	EN_US string
	PT_BR string
	FR_FR string
	ES_ES string
	ZH_CN string
}

type TELEMETRY struct {
	END_POINT    string
	SERVICE_NAME string
}

type SECURITY struct {
	RESET_PASSWORD_EXPIRATION_TIME string
	EMAIL_CHANGE_EXPIRATION_TIME   string
	VERIFY_EMAIL_EXPIRATION_TIME   string
	MAX_CHANGE_EMAIL_ATTEMPTS      string
}

type EMAIL_SERVICE struct {
	API_KEY string
	VERIFY  string
}

type SERVER struct {
	PORT string
}

var DB_POSTGRES = DB{
	DB_HOST:     os.Getenv("DB_HOST"),
	DB_USER:     os.Getenv("DB_USER"),
	DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	DB_PORT:     os.Getenv("DB_PORT"),
	DB_NAME:     os.Getenv("DB_NAME"),
}

var SECRETS_VAR = SECRETS{
	JWT_SECRET:         os.Getenv("JWT_SECRET"),
	JWT_ACCESS_SECRET:  os.Getenv("JWT_ACCESS_SECRET"),
	JWT_REFRESH_SECRET: os.Getenv("JWT_REFRESH_SECRET"),
}

var FRONT_END_URL_VAR = FRONT_END_URL{
	FRONT_END_URL_DEV:  os.Getenv("FRONT_END_URL_DEV"),
	FRONT_END_URL_PROD: os.Getenv("FRONT_END_URL_PROD"),
}

var GOOGLE_VAR = GOOGLE{
	IMAGE_BUCKET_NAME: os.Getenv("IMAGE_BUCKET_NAME"),
	URL_BUCKET_NAME:   os.Getenv("URL_BUCKET_NAME"),
}

var AVAILABLE_LANGUAGES_VAR = AVAILABLE_LANGUAGES{
	EN_US: "en-US",
	PT_BR: "pt-BR",
	FR_FR: "fr-FR",
	ES_ES: "es-ES",
	ZH_CN: "zh-CN",
}

var SECURITY_VAR = SECURITY{
	RESET_PASSWORD_EXPIRATION_TIME: os.Getenv("RESET_PASSWORD_EXPIRATION_TIME"),
	EMAIL_CHANGE_EXPIRATION_TIME:   os.Getenv("EMAIL_CHANGE_EXPIRATION_TIME"),
	VERIFY_EMAIL_EXPIRATION_TIME:   os.Getenv("VERIFY_EMAIL_EXPIRATION_TIME"),
	MAX_CHANGE_EMAIL_ATTEMPTS:      os.Getenv("MAX_CHANGE_EMAIL_ATTEMPTS"),
}

var EMAIL_SERVICE_VAR = EMAIL_SERVICE{
	API_KEY: os.Getenv("EMAIL_SERVICE_API_KEY"),
	VERIFY:  os.Getenv("EMAIL_VERIFY"),
}

var SERVER_VAR = SERVER{
	PORT: os.Getenv("SERVER_PORT"),
}

func GetFrontendURLs() (string, string) {
	return os.Getenv("FRONT_END_URL_DEV"), os.Getenv("FRONT_END_URL_PROD")
}

func GetBucketNames() (string, string) {
	return os.Getenv("IMAGE_BUCKET_NAME"), os.Getenv("URL_BUCKET_NAME")
}

func GetEmailChangeExpirationTime() string {
	return os.Getenv("EMAIL_CHANGE_EXPIRATION_TIME")
}

func GetMaxChangeEmailAttempts() string {
	return os.Getenv("MAX_CHANGE_EMAIL_ATTEMPTS")
}

func GetEmailServiceAPIKey() string {
	return os.Getenv("EMAIL_SERVICE_API_KEY")
}

func GetEmailServiceVerify() string {
	return os.Getenv("EMAIL_VERIFY")
}

func GetServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return ":8080"
	}
	if port[0] != ':' {
		return ":" + port
	}
	return port
}

func GetSwaggerHost() string {
	return os.Getenv("SWAGGER_HOST")
}
