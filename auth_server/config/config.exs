import Config

config :auth_server, AuthServer.Repo,
  database: "auth_server_repo",
  username: "gerri",
  password: "H@lunk3nm4nn",
  hostname: "localhost"

config :auth_server,
  ecto_repos: [AuthServer.Repo]
