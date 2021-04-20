package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// BannedIdent struct which all information about a banned identifier:
// - fully qualified import/function name,
// - message with additional information,
// - a list of exemptions.
type BannedIdent struct {
	Name       string      `json:"name"`
	Msg        string      `json:"msg"`
	Exemptions []Exemption `json:"exemptions"`
}

// BannedIdents is a map of identifier names to a list of BannedIdent
// entries that define additional information.
type BannedIdents map[string][]BannedIdent

// Exemption struct which contains a justification and a path to a directory
// that should be exempted from the check.
type Exemption struct {
	Justification string `json:"justification"`
	AllowedDir    string `json:"allowedDir"`
}

// ReadBannedImports reads banned imports from all config files
// and concatenates them into one object.
func ReadBannedImports(files []string) (BannedIdents, error) {
	imports := make(BannedIdents)

	for _, file := range files {
		config, err := unmarshalCfg(file)
		if err != nil {
			return nil, err
		}

		for _, i := range config.Imports {
			imports[i.Name] = append(imports[i.Name], i)
		}
	}

	return imports, nil
}

// ReadBannedFunctions reads banned function calls from all config files
// and concatenates them into a map.
func ReadBannedFunctions(files []string) (BannedIdents, error) {
	fns := make(BannedIdents)

	for _, file := range files {
		config, err := unmarshalCfg(file)
		if err != nil {
			return nil, err
		}

		for _, fn := range config.Functions {
			fns[fn.Name] = append(fns[fn.Name], fn)
		}
	}

	return fns, nil
}

// config struct which contains an array of banned imports and function calls.
type config struct {
	Imports   []BannedIdent `json:"imports"`
	Functions []BannedIdent `json:"functions"`
}

// unmarshalCfg reads JSON object from a file and converts it to a config struct.
func unmarshalCfg(file string) (*config, error) {
	if !fileExists(file) {
		return nil, errors.New("file does not exist or is a directory")
	}

	cfg, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer cfg.Close()

	var config config
	bytes, _ := ioutil.ReadAll(cfg)
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}