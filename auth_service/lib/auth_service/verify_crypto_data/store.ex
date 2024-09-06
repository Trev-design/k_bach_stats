defmodule AuthService.VerifyCryptoData.Store do
  use GenServer

  alias AuthService.VerifyCryptoData.{PurgeHelper, StoreHelpers, HandlerFunctions}

  defstruct key: nil, ttl: 0

  def start_link(props), do: GenServer.start_link(__MODULE__, props, name: __MODULE__)

  def init([key: key, ttl: ttl]) do
    HandlerFunctions.create_store()
    {:ok, %__MODULE__{key: key, ttl: ttl}, {:continue, :prepare_purge}}
  end

  def encrypt(id, value, pid), do: GenServer.cast(__MODULE__, {:encrypt, id, value, pid})
  def decrypt(id, cypher, pid), do: GenServer.cast(__MODULE__, {:decrypt, id, cypher, pid})

  def handle_cast({:encrypt, id, value, pid}, state) do
    result = Task.await(StoreHelpers.encrypt(id, state.key, value))
    send(pid, result)
    {:noreply, state}
  end

  def handle_cast({:decrypt, id, cypher, pid}, state) do
    result = Task.await(StoreHelpers.decrypt(id, state.key, cypher))
    send(pid, result)
    {:noreply, state}
  end

  def handle_continue(:prepare_purge, state) do
    Process.send_after(self(), :start_purge, state.ttl)
    {:noreply, state}
  end

  def handle_info(:start_purge, state) do
    PurgeHelper.prepare_purge()
    {:noreply, state, {:continue, :prepare_purge}}
  end
end
