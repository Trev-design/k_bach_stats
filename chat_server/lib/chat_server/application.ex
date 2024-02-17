defmodule ChatServer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      {
        Plug.Cowboy,
        scheme: :http,
        plug: ChatServer.Router,
        options: [
          dispatch: dispatch(),
          port: 4001
        ]
      },

      {Registry, keys: :duplicate, name: Registry.ChatServer},
      {Task.Supervisor, name: ChatServer.TaskSupervisor}
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: ChatServer.Supervisor]
    Supervisor.start_link(children, opts)
  end

  defp dispatch() do
    [
      {:_,
        [
          {"/ws/[...]", ChatServer.SocketHandler, []},
          {:_, Plug.Cowboy.Handler, {ChatServer.Router, []}}
        ]
      }
    ]
  end
end
