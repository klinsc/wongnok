package config

type Keycloak struct {
	ClientID     string `env:"KEYCLOAK_CLIENT_ID" envDefault:"wongnok"`
	ClientSecret string `env:"KEYCLOAK_CLIENT_SECRET"`
	RedirectURL  string `env:"KEYCLOAK_REDIRECT_URL" envDefault:"http://localhost:8080/api/v1/callback"`
	Realm        string `env:"KEYCLOAK_REALM" envDefault:"pea-devpool-2025"`
	URL          string `env:"KEYCLOAK_URL" envDefault:"https://sso-dev.odd.works"`
}
