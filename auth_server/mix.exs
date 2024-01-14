defmodule AuthServer.MixProject do
  use Mix.Project

  def project do
    [
      app: :auth_server,
      version: "0.1.0",
      elixir: "~> 1.16",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {AuthServer.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:bandit, "~> 1.1"},
      {:postgrex, "~> 0.17.4"},
      {:ecto, "~> 3.11"},
      {:ecto_sql, "~> 3.11"},
      {:jason, "~> 1.4"},
      {:comeonin, "~> 5.4"},
      {:argon2_elixir, "~> 4.0"},
      {:joken, "~> 2.6"},
      {:uuid, "~> 1.1"}
    ]
  end
end
