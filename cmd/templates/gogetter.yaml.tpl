# Project Configuration

module: "{{.Module}}"
dir: "{{.Dir}}"

# MANDATORY MODULES
logger:
  include: {{.IncludeLogger}}
  type: "{{.Logger}}"

config:
  include: {{.IncludeConfig}}
  type: "{{.ConfigManager}}"

http_server:
  include: {{.IncludeHTTPServer}}
  type: "{{.Framework}}"

# OPTIONAL MODULES
jwt_middleware:
  include: {{.IncludeJWTMiddleware}}
  type: "{{.JWTMiddleware}}"

db_connector:
  include: {{.IncludeDBConnector}}
  type: "{{.DBConnector}}"

mailer:
  include: {{.IncludeMailer}}
  type: "{{.Mailer}}"

git:
  set_up: {{.SetUpGit}}

docker:
  set_up_file: {{.SetUpDockerFile}}
  set_up_compose: {{.SetUpDockerCompose}}