package model

type Config struct {
	Database DatabaseConfig `json:"database"`
	Jwt      JwtConfig      `json:"jwt"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type JwtConfig struct {
	AccessName    string `json:"access_name"`
	AccessSecret  string `json:"access_secret"`
	AccessExp     int    `json:"access_exp"`
	RefreshName   string `json:"refresh_name"`
	RefreshSecret string `json:"refresh_secret"`
	RefreshExp    int    `json:"refresh_exp"`
}
