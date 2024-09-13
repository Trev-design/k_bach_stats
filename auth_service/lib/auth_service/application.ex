defmodule AuthService.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      AuthServiceWeb.Telemetry,
      AuthService.Repo,
      {DNSCluster, query: Application.get_env(:auth_service, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: AuthService.PubSub},
      AuthServiceWeb.Endpoint,

      {AuthService.VerifyCryptoData.Store, verify_storage_props()},
      {Poolex, purge_pool()},
      {Poolex, crypto_pool()},
      {Task.Supervisor, name: PurgeHelper.Supervisor},

      {Redix, redix_spec()},

      {AuthService.Rabbitmq.RabbitInstance, []}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: AuthService.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    AuthServiceWeb.Endpoint.config_change(changed, removed)
    :ok
  end

  defp verify_storage_props() do
    [
      key: :crypto.strong_rand_bytes(32),
      ttl: 10000
    ]
  end

  defp purge_pool() do
    [
      pool_id: :purge,
      worker_module: AuthService.VerifyCryptoData.Purger,
      workers_count: 4,
      max_overflow: 1
    ]
  end

  defp crypto_pool() do
    [
      pool_id: :crypto,
      worker_module: AuthService.VerifyCryptoData.CryptoHandler,
      workers_count: 8,
      max_overflow: 2
    ]
  end

  defp redix_spec() do
    [
      host: "localhost",
      port: 6379,
      database: 0,
      name: :verify_session_store
    ]
  end
end
