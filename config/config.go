package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Load loads the configuration from the given file path into the dst struct.
func Load(dst any, path string) error {
	viper.AddConfigPath(fileDir(path))
	viper.SetConfigName(filename(path))
	viper.SetConfigType(fileType(path))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed viper.ReadInConfig: %w", err)
	}

	err = viper.Unmarshal(dst)
	if err != nil {
		return fmt.Errorf("failed viper.Unmarshal: %w", err)
	}

	return nil
}

// fileType extracts the file extension from the path.
func fileType(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

// filename extracts the filename without the extension.
func filename(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// fileDir extracts the directory path.
func fileDir(path string) string {
	return filepath.Dir(path)
}
