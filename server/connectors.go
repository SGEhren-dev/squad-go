package server

type DatabaseConnector struct {
	DatabaseName string `json:"databaseName"`
	Hostname     string `json:"hostname"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Port         int    `json:"port"`
}
