defmodule AuthService.VerifyCryptoData.Access do
  alias AuthService.VerifyCryptoData.Store

  def encrypted(id, value), do: compute(fn -> Store.encrypt(id, value, self()) end)

  def decrypted(id, cypher), do: compute(fn -> Store.decrypt(id, cypher, self()) end)

  defp compute(func) do
    func.()
    receive do
      {:ok, _} = result -> result
      _invalid          -> :error
    end
  end
end
