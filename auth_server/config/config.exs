import Config

config :auth_server, AuthServer.Repo,
  database: "auth_server_repo",
  username: "IAmTheUser",
  password: "ThisIsMyPassword",
  hostname: "localhost"

config :auth_server,
  ecto_repos: [AuthServer.Repo]
