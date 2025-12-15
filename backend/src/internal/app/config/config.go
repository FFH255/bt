package config

type Config struct {
	Server     Server     `yaml:"server"`
	Logger     Logger     `yaml:"logger"`
	CORS       CORS       `yaml:"cors"`
	Redis      Redis      `yaml:"redis"`
	Swagger    Swagger    `yaml:"swagger"`
	Cookie     Cookie     `yaml:"cookie"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Auth       Auth       `yaml:"auth"`
	Postgres   Postgres   `yaml:"postgres"`
	Profile    Profile    `yaml:"profile"`
	Statistics Statistics `yaml:"statistics"`
	Antifroad  Antifroad  `yaml:"antifroad"`
	Languages  []string   `yaml:"languages"`
}

type Server struct {
	Address         string `yaml:"address"`
	ReadTimeout     string `yaml:"read_timeout"`
	WriteTimeout    string `yaml:"write_timeout"`
	IdleTimeout     string `yaml:"idle_timeout"`
	ShutdownTimeout string `yaml:"shutdown_timeout"`
}

type Logger struct {
	Enabled bool   `yaml:"enabled"`  // Нужно ли писать логи?
	UseFile bool   `yaml:"use_file"` // Писать логи в файл иил в stdout?
	Path    string `yaml:"path"`     // Путь к файлу с логами (если use_file = true)
}

type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowAllOrigins  bool     `yaml:"allow_all_origins"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

type Postgres struct {
	Connection string `yaml:"connection"`
	Migrations string `yaml:"migrations"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type Providers struct {
	Google Provider `yaml:"google"`
	Github Provider `yaml:"github"`
}

type Provider struct {
	ClientKey   string `yaml:"client_key"`
	Secret      string `yaml:"secret"`
	CallbackURL string `yaml:"callback_url"`
}

type Auth struct {
	JWTSecret               string    `yaml:"jwt_secret"`
	Providers               Providers `yaml:"providers"`
	LoggedInRedirectURL     string    `yaml:"logged_in_redirect_url"`
	RegistrationRedirectURL string    `yaml:"registration_redirect_url"`
	ErrorRedirectURL        string    `yaml:"error_redirect_url"`
}

type Swagger struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type Cookie struct {
	Secure bool   `yaml:"secure"`
	Key    string `yaml:"key"`
}

type Scheduler struct {
	DeleteExpiredSessionsInterval string `yaml:"delete_expired_sessions_interval"`
}

type Profile struct {
	Expiration string `yaml:"expiration"`
	UseRedis   bool   `yaml:"use_redis"`
}

type Statistics struct {
	PBExpirationTime string `yaml:"pb_expiration_time"`
}

type Antifroad struct {
	IsDisabled       bool   `yaml:"is_disabled"`       // Нужно ли выключить модуль антифрода?
	MaxKeys          int    `yaml:"max_keys"`          // Максимальное кол-во одновременно существующих в базе ключей
	RotationInterval string `yaml:"rotation_interval"` // Интервал ротации ключей в базе
	Password         string `yaml:"password"`          // Пароль для авторизации в модуле антифрода внутри контура
}
