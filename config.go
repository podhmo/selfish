package selfish

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

// Config is mapping object for application config
type Config struct {
	DefaultAlias string `json:"default_alias"`
	AccessToken  string `json:"access_token"`
	HistFile     string `json:"hist_file"`
}

var (
	defaultConfigDir string
	defaultAlias     string
)

func init() {
	defaultConfigDir = path.Join(os.Getenv("HOME"), ".selfish")
	defaultAlias = "head"
	// fmt.Printf("history: %q, alias: %q\n", defaultHistFile, defaultAlias)
}

// Dirs returns config-directory's candidates
func Dirs() []string {
	candidates := []string{".", os.Getenv("HOME")}
	return dirs(candidates)
}

func dirs(candidates []string) []string {
	var paths []string
	for _, d := range candidates {
		paths = append(paths, path.Join(d, ".selfish"))
	}
	return append(paths, defaultConfigDir)
}

// GetConfigDir returns a path of config directory
func GetConfigDir() (string, error) {
	for _, path := range Dirs() {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return defaultConfigDir, errors.Errorf("config directory is not found. (default is ~/.selfish)")
}

// LoadConfig loads configuration file, if configuration file is not existed, then return default config.
func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, errors.Errorf("%q is not found (dir)", configDir)
	}
	return loadConfig(configDir)
}

func loadConfig(d string) (*Config, error) {
	filename := path.Join(d, "config.json")
	if _, err := os.Stat(filename); err != nil {
		return nil, errors.Wrap(err, "stat")
	}
	fp, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, errors.Wrap(err, "read all")
	}

	config := Config{}
	json.Unmarshal(data, &config)

	if config.HistFile == "" {
		config.HistFile = path.Join(d, "selfish.history")
	}
	if config.DefaultAlias == "" {
		config.DefaultAlias = defaultAlias
	}
	return &config, nil
}
