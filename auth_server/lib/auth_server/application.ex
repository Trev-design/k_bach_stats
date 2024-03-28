defmodule AuthServer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      AuthServer.VerifyStore,
      {
        Bandit,
        scheme: :http,
        plug: AuthServer.Router,
        port: 4000
      },
      AuthServer.Repo,
      {Task.Supervisor, name: MailRequest.Supervisor}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: AuthServer.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
