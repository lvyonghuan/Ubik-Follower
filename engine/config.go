package engine

const configPath = "./conf/config.json"
const testConfigPath = "../conf/config.json"

type Config struct {
	Port       string `json:"port"`
	PluginPath string `json:"plugin_path"`
	Log        Log    `json:"log"`

	LeaderUrl string `json:"leader_url"` //Leader url, if null will broadcast to find leader
}

type Log struct {
	Level       int    `json:"level"`     //Log level
	IsSave      bool   `json:"is_save"`   //Whether to save logs
	LogSavePath string `json:"save_path"` //The path where the logs are saved
}
