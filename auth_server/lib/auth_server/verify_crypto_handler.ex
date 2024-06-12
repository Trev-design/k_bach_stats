defmodule AuthServer.VerifyCryptoHandler do
  use GenServer

  def start_link(_), do: GenServer.start_link(__MODULE__, :ok, name: __MODULE__)

  @impl GenServer
  def init(_) do
    :mnesia.create_schema([node()])
    :mnesia.start()
    create_table()

    {:ok, :crypto.strong_rand_bytes(32)}
  end

  def encrypt(pid, id, plain), do: GenServer.cast(__MODULE__, {:encrypt, pid, id, plain})
  def decrypt(pid, id, cypher), do: GenServer.cast(__MODULE__, {:decrypt, pid, id, cypher})

  def purge(), do: GenServer.cast(__MODULE__, :purge)

  @impl GenServer
  def handle_cast({:encrypt, pid, id, plain}, key) do
    send(pid, encrypt_response(id, plain, key))
    {:noreply, key}
  end

  def handle_cast({:decrypt, pid, id, cypher}, key) do
    send(pid, decrypt_response(id, cypher, key))
    {:noreply, key}
  end

  @impl GenServer
  def handle_cast(:purge, state) do
    do_purge()
    {:noreply, state}
  end

  defp encrypt_response(id, plain, key) do
    iv = :crypto.strong_rand_bytes(16)
    cypher = :crypto.crypto_one_time(:aes_256_ofb, key, iv, plain, true)
    func = fn -> :mnesia.write({VerifyCryptoData, id, iv,:erlang.system_time(:millisecond) + (150 * 60 * 1000)}) end
    case :mnesia.transaction(func) do
      {:atomic, :ok}     -> {:ok, cypher}
      {:aborted, reason} -> {:error, reason}
    end
  end

  defp decrypt_response(id, cypher, key) do
    func = fn ->
      case :mnesia.read({VerifyCryptoData, id}) do
        [{VerifyCryptoData, ^id, _iv, _created_at} = data] -> data

        _invalid ->
          :mnesia.abort("could not find data")
      end
    end

    case :mnesia.transaction(func) do
      {:atomic, {VerifyCryptoData, ^id, iv, _created_at}} ->
        plain = :crypto.crypto_one_time(:aes_256_ofb, key, iv, cypher, false)
        :mnesia.transaction(fn -> :mnesia.delete({VerifyCryptoData, id}) end)
        {:ok, plain}

      {:aborted, "could not find data" = reason} -> {:error, reason}

      {:aborted, reason} -> {:error, elem(reason, 0) |> Atom.to_string()}
    end
  end

  defp do_purge() do
    :mnesia.sync_transaction(
      fn ->
        current_time = :erlang.system_time(:millisecond)
        case :mnesia.select(VerifyCryptoData, [{{:"$1", :"$2", :"$3"}, [:<=, :"$3", current_time], [:"$1"]}]) do
          []    -> :mnesia.abort("no selected terms")
          terms ->
            Enum.each(terms, fn id ->
              :mnesia.delete({VerifyCryptoData, id})
            end)
        end
      end)
  end

  defp create_table() do
    :mnesia.table_info(VerifyCryptoData, :all)
  catch
    :exit, _ ->
      :mnesia.create_table(
        VerifyCryptoData,
        [
          type: :set,
          attributes: [:id, :iv, :ttl],
          disc_copies: [node()],
        ])
  end
end
