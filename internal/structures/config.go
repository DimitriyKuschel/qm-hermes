package structures

import "time"

type Template struct {
	TplDir    string `yaml:"tplDir" validate:"required"`
	PublicDir string `yaml:"publicDir" validate:"required"`
}

type Server struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required|uint|min:1"`
}

type Persistence struct {
	FilePath     string        `yaml:"filePath" validate:"required|unixPath"`
	SaveInterval time.Duration `yaml:"saveInterval" validate:"required|min:1"`
}

type LoggerConfig struct {
	Level string `yaml:"level" validate:"required|in:trace,debug,info,warn,error,fatal,panic"`
	Mode  uint32 `yaml:"mode" validate:"required|uint"`
	Dir   string `yaml:"dir" validate:"required|unixPath"`
}

type DashboardAuthentication struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Required bool   `yaml:"required" validate:"bool"`
	Secret   string `yaml:"secret" validate:"required"`
}

type Config struct {
	AppName                 string
	Debug                   bool
	Path                    string
	WebServer               Server                  `yaml:"webServer"`
	TcpServer               Server                  `yaml:"tcpServer"`
	Persistence             Persistence             `yaml:"persistence"`
	Logger                  LoggerConfig            `yaml:"logger"`
	Template                Template                `yaml:"template"`
	DashboardAuthentication DashboardAuthentication `yaml:"dashboardAuthentication"`
}
