defmodule AuthService.VerifyCryptoData.StoreHelpers do
  def encrypt(id, key, value) do
    make_task(fn pid -> GenServer.call(pid, {:encrypt, id, key, value})end)
  end

  def decrypt(id, key, cypher) do
    make_task(fn pid -> GenServer.call(pid, {:decrypt, id, key, cypher}) end)
  end

  defp make_task(fun) do
    Task.async(fn ->
      Poolex.run(
        :crypto,
        fun,
        checkout_timeout: 10000
      )
    end)
  end
end
