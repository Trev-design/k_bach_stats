import Config

config :auth_server, AuthServer.Repo,
  database: "auth_server_repo",
  username: "gerri",
  password: "H@lunk3nm4nn",
  hostname: "localhost"

config :auth_server, AuthServer.JobHandler.Repo,
  database: "auth_server_jobs_repo",
  username: "gerri",
  password: "H@lunk3nm4nn",
  hostname: "localhost"

config :auth_server,
  ecto_repos: [AuthServer.Repo, AuthServer.JobHandler.Repo]

config :auth_server, Oban,
  repo: AuthServer.JobHandler.Repo,
  plugins: [Oban.Plugins.Pruner],
  queues: [events: [limit: 2, paused: false]]
