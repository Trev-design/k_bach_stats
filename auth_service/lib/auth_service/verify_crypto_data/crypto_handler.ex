defmodule AuthService.VerifyCryptoData.CryptoHandler do
  alias AuthService.VerifyCryptoData.HandlerFunctions
  use GenServer


  def start_link(), do: GenServer.start_link(__MODULE__, nil)

  @impl GenServer
  def init(nil), do: {:ok, nil}

  @impl GenServer
  def handle_call({:encrypt, id, key, value}, _from, state) do
    {:reply, HandlerFunctions.encrypt_transaction(id, key, value), state}
  end

  @impl GenServer
  def handle_call({:decrypt, id, key, cypher}, _from, state) do
    {:reply, HandlerFunctions.decrypt_transaction(id, key, cypher), state}
  end
end
