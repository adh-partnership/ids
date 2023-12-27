package config

type Config struct {
	Cache    Cache    `json:"cache"`
	Database Database `json:"database"`
	Facility Facility `json:"facility"`
	OAuth    OAuth    `json:"oauth"`
	Server   Server   `json:"server"`
	Session  Session  `json:"session"`
}

type Cache struct {
	Driver            string           `json:"driver"`
	Host              string           `json:"host"`
	Port              int              `json:"port"`
	Password          string           `json:"password"`
	DB                int              `json:"db"`
	DefaultExpiration *CacheExpiration `json:"default_expiration"`
}

type CacheExpiration struct {
	Airports int `json:"airports"`
	Charts   int `json:"charts"`
	PIREPs   int `json:"pireps"`
}

type Database struct {
	Driver       string `json:"driver"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	AutoMigrate  bool   `json:"auto_migrate"`
	CACert       string `json:"ca_cert"`
}

type Facility struct {
	Identifier string      `json:"identifier"`
	Name       string      `json:"name"`
	ADH        FacilityADH `json:"adh"`
}

type FacilityADH struct {
	APIBase  string `json:"api_base"`
	Rostered bool   `json:"rostered"`
}

type OAuth struct {
	Provider     string         `json:"provider"`
	BaseURL      string         `json:"base_url"`
	ClientID     string         `json:"client_id"`
	ClientSecret string         `json:"client_secret"`
	MyBaseURL    string         `json:"my_base_url"`
	Endpoints    OAuthEndpoints `json:"endpoints"`
}

type OAuthEndpoints struct {
	Authorization string `json:"authorize"`
	Token         string `json:"token"`
	Userinfo      string `json:"userinfo"`
}

type Server struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Mode string `json:"mode"`
}

type Session struct {
	HashKey  string `json:"hash_key"`
	BlockKey string `json:"block_key"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	MaxAge   int    `json:"max_age"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}
