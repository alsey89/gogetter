package cmd

type ProjectConfig struct {
	Module string
	Dir    string
	// MANDATORY MODULES
	// -----logger-----
	IncludeLogger bool   `json:"include_logger" yaml:"include_logger"`
	Logger        string `json:"logger" yaml:"logger"`
	// -----config-----
	IncludeConfig bool   `json:"include_config" yaml:"include_config"`
	ConfigManager string `json:"config_manager" yaml:"config_manager"`
	// -----http server-----
	IncludeHTTPServer bool   `json:"include_http_server" yaml:"include_http_server"`
	Framework         string `json:"framework" yaml:"framework"`

	// OPTIONAL MODULES
	// -----jwt middleware-----
	IncludeJWTMiddleware bool   `json:"include_jwt_middleware" yaml:"include_jwt_middleware"`
	JWTMiddleware        string `json:"jwt_middleware" yaml:"jwt_middleware"`
	// -----db connector-----
	IncludeDBConnector bool   `json:"include_db_connector" yaml:"include_db_connector"`
	DBConnector        string `json:"db_connector" yaml:"db_connector"`
	// -----mailer-----
	IncludeMailer bool   `json:"include_mailer" yaml:"include_mailer"`
	Mailer        string `json:"mailer" yaml:"mailer"`
	// -----git-----
	SetUpGit bool `json:"set_up_git" yaml:"set_up_git"`
	// -----docker-----
	SetUpDockerFile    bool `json:"set_up_docker_file" yaml:"set_up_docker_file"`
	SetUpDockerCompose bool `json:"set_up_docker_compose" yaml:"set_up_docker_compose"`
}
