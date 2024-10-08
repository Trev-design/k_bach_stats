# This file is responsible for configuring your application
# and its dependencies with the aid of the Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
import Config

config :auth_service,
  ecto_repos: [AuthService.Repo],
  generators: [timestamp_type: :utc_datetime, binary_id: true]

# Configures the endpoint
config :auth_service, AuthServiceWeb.Endpoint,
  url: [host: "localhost"],
  adapter: Bandit.PhoenixAdapter,
  render_errors: [
    formats: [json: AuthServiceWeb.ErrorJSON],
    layout: false
  ],
  pubsub_server: AuthService.PubSub,
  live_view: [signing_salt: "vrTdxjnq"]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

config  :auth_service, Oban,
  repo: AuthService.Repo,
  queues: [verify_user: 2, authorize_user: 2, enroll_user: 2]

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
