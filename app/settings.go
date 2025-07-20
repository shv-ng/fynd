package app

type Settings struct {
	RootPath       string   `yaml:"root_path"`
	DBPath         string   `yaml:"db_path"`
	Top            int      `yaml:"top"`
	MaxConcurrency uint     `yaml:"max_concurrency"`
	IncludeDirs    []string `yaml:"include_dirs"`
	ExcludeDirs    []string `yaml:"exclude_dirs"`
	IncludeHidden  bool     `yaml:"include_hidden"`
}
