package engine

const (
	confPath     = "./conf"
	testConfPath = "../conf"
)

type Config struct {
	Port              string `json:"port"`
	IP                string `json:"ip"`
	PluginPath        string `json:"plugin_path"`
	HeartBeatInterval int    `json:"heartbeat_interval"` //Heartbeat interval in seconds
	Log               Log    `json:"log"`

	LeaderAddr string `json:"leader_addr"` //Leader url, if null will broadcast to find leader
}

type Log struct {
	Level       int    `json:"level"`     //Log level
	IsSave      bool   `json:"is_save"`   //Whether to save logs
	LogSavePath string `json:"save_path"` //The path where the logs are saved
}
