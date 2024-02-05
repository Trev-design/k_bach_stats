defmodule AuthServer.VerifyStore do
  use GenServer

  def start_link(_opts), do: GenServer.start_link(__MODULE__, %{}, name: __MODULE__)

  def init(_init_arg) do
    :ets.new(:verify_session_table, [:named_table, read_concurrency: true])
    {:ok, :ok}
  end
end
