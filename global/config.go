package global

// ConfigFile - is the config file path
const ConfigFile string = "./config/config.yml"

// Config - the config file structure represented as struct
type Config struct {
	Server string `yaml:"server"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

// WriteConfig - is w
func (c *Config) WriteConfig() {

}
