package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/shv-ng/fynd/app"
)

func LoadYAMLSettings(path string) (app.Settings, error) {
	var s app.Settings
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("No existing config file, generated one in", path)
		generateFile(path)
	}
	err = yaml.Unmarshal(data, &s)
	return s, err
}

func generateFile(path string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get homedir: %v", err)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("failed to get cachedir: %v", err)
	}
	c := fmt.Sprintf(`root_path: %v
db_path: %v
top: 15
max_concurrency: 150
include_hidden: false
include_dirs:
exclude_dirs:
  # Version control and IDE
  - .git
  - .svn
  - .hg
  - .idea
  - .vscode
  - .metadata

  # Python
  - __pycache__
  - .venv
  - venv
  - .env
  - env
  - .mypy_cache
  - .pytest_cache
  - .tox
  - .eggs
  - .coverage
  - htmlcov
  - .cache
  - pip-wheel-metadata
  - site

  # JavaScript / TypeScript
  - node_modules
  - .yarn
  - .yarn-cache
  - .pnp
  - .parcel-cache
  - .next
  - .nuxt
  - .output
  - .angular
  - jspm_packages
  - bower_components
  - coverage
  - dist
  - build
  - .turbo
  - .expo

  # Rust
  - target
  - .cargo

  # Go
  - bin
  - pkg
  - vendor

  # Dart / Flutter
  - .dart_tool
  - .flutter-plugins
  - .flutter-plugins-dependencies
  - .packages
  - .android
  - .ios
  - ios
  - android
  - linux
  - macos
  - windows
  - build

  # Java / Kotlin
  - .gradle
  - .settings
  - out
  - libs

  # C / C++ / CMake
  - CMakeFiles

  # Misc
  - tmp
  - temp
  - logs
  - log
  - backup
  - backups`,
		home,
		filepath.Join(cache, "fynd", "files.db"),
	)
	os.Mkdir(filepath.Dir(path), 0o700)
	os.Mkdir(filepath.Join(cache, "fynd"), 0o700)
	if err := os.WriteFile(path, []byte(c), 0o600); err != nil {
		log.Fatalf("Failed to write content: %v\n", err)
	}
}
