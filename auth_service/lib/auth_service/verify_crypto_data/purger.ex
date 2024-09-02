defmodule AuthService.VerifyCryptoData.Purger do
  alias AuthService.VerifyCryptoData.HandlerFunctions
  use GenServer

  def start_link(), do: GenServer.start(__MODULE__, nil)

  @impl GenServer
  def init(nil), do: {:ok, nil}

  @impl GenServer
  def handle_call({:handle_purge, terms}, _from, state) do
    result = HandlerFunctions.purge_transaction(terms)
    Process.sleep(5000)
    {:reply, result, state}
  end
end
