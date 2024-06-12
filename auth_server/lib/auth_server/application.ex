defmodule AuthServer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      AuthServerWeb.Telemetry,
      AuthServer.VerifyCryptoHandler,
      AuthServer.Repo,
      {
        Redix,
        host: "localhost",
        port: 6379,
        database: 0,
        name: :verify_store
      },
      {DNSCluster, query: Application.get_env(:auth_server, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: AuthServer.PubSub},
      {Finch, name: EmailRequest},
      # Start a worker by calling: AuthServer.Worker.start_link(arg)
      # {AuthServer.Worker, arg},
      # Start to serve requests, typically the last entry
      AuthServerWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: AuthServer.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    AuthServerWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
