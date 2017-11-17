package properties

type Properties struct {
	originPath string
	isLoaded   bool

	System  SystemDefinition   `json:"system" mapstructure:"system"`
	Modules []ModuleDefinition `json:"modules" mapstructure:"modules"`
}

type ModuleDefinition struct {
	Name string      `json:"name" mapstructure:"name"`
	Path string      `json:"path" mapstructure:"path"`
	Conf interface{} `json:"conf" mapstructure:"conf"`
}

type SystemDefinition struct {
	LinkID     string `json:"link_id" mapstructure:"conf"`
	ServerIP   string `json:"ip" mapstructure:"ip"`
	ServerPort int    `json:"port" mapstructure:"port"`
	ServerSSL  bool   `json:"ssl" mapstructure:"ssl"`
}

func (p *Properties) IsLoaded() bool {
	return p.isLoaded
}
