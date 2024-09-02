defmodule AuthService.VerifyCryptoData.HandlerFunctions do
  def create_store() do
    :mnesia.start()
    :mnesia.create_schema([node()])
    create_table([node()])
  end

  def encrypt_transaction(id, key, value) do
    iv = :crypto.strong_rand_bytes(16)
    cypher = :crypto.crypto_one_time(:aes_256_ofb, key, iv, value, true)
    ttl = :erlang.system_time(:second) + 60 * 120
    transaction = fn ->:mnesia.write({VerifyData, id, iv, ttl, false}) end

    case :mnesia.activity(:transaction, transaction, [], :mnesia_frag) do
      :ok -> cypher
      _   -> :something_went_wrong
    end
  end

  def decrypt_transaction(id, key, cypher) do
    transaction = fn ->
      case :mnesia.read(VerifyData, id) do
        [data] -> get_data(data)
        []     -> :mnesia.abort(:could_not_find_data)
      end
    end

    case :mnesia.activity(:transaction, transaction, [], :mnesia_frag) do
      {:atomic, {_, _, iv, _, _}} ->
        :crypto.crypto_one_time(:aes_256_0fb, key, iv, cypher, false)

      _invalid -> :something_went_wrong
    end
  end

  def get_deletables_transaction() do
    transaction = fn ->
      current_time = :erlang.system_time(:second)

      case :mnesia.select(
        VerifyData,
        [{{:"$1", :"$2", :_, :"$3", :"$4"},
        [{:or, {:<, :"$3", current_time}, {:==, :"$4", true}}],
        [{:"$1", :"$2"}]}])
      do
        []    -> nil
        terms -> terms
      end
    end

    case :mnesia.activity(:transaction, transaction, [], :mnesia_frag) do
      {:atomic, terms} -> terms
      {:aborted, _}    -> nil
    end
  end

  def purge_transaction(terms) do
    transaction = fn ->
      Enum.each(terms, fn {VerifyData, id} ->
        :mnesia.delete(VerifyData, id, :write)
      end)
    end

    :mnesia.activity(:transaction, transaction, [], :mnesia_frag)
  end

  defp get_data({VerifyData, id, iv, ttl, _}) do
    new = {VerifyData, id, iv, ttl, true}
    :mnesia.write(new)
    new
  end

  defp create_table(nodes) do
    :mnesia.table_info(Verifydata, :all)

  catch
    :exit, _ ->
      :mnesia.create_table(
        VerifyData,
        attributes: [:id, :iv, :ttl, :selected],
        type: :set,
        ram_copies: nodes,
        frag_properties: [
          n_fragments: 10,
          node_pool: nodes,
          n_ram_copies: 1
        ]
      )
  end
end
