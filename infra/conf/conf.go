package conf

type Conf struct {
	ProductionMode bool          `json:"production"`
	Logger         *Logger       `json:"logger"`
	Web            *Web          `json:"web"`
	Database       *DatabaseConf `json:"database"`
}

type Logger struct {
	Level string `json:"level"`
}

type Web struct {
	Port int `json:"port"`
}

type DatabaseConf struct {
	Driver   string `json:"driver"`   // 数据库驱动类型
	Host     string `json:"host"`     // 数据库主机地址
	Port     int    `json:"port"`     // 数据库端口
	Database string `json:"database"` // 数据库名称
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	DSN      string `json:"dsn"`      // 完整的数据源名称，如果提供则优先使用
}
