defmodule AuthServer.VerifyHandler do
  alias AuthServer.VerifyCryptoHandler

  require Logger

  def encrypt_verification_code(id, verification) do
    VerifyCryptoHandler.encrypt(self(), id, Jason.encode!(%{verification: verification}))
    receive do
      {:ok, _verification} = result -> result
      invalid                       -> invalid
    end
  end

  def decrypt_verification_code(id, verification) do
    VerifyCryptoHandler.decrypt(self(), id, verification)
    receive do
      {:ok, verify} ->
        case Jason.decode(verify) do
          {:ok, %{"verification" => verification}} ->
            {:ok, verification}

          {:error, reason} ->
            IO.inspect(reason)
            {:error, reason.data()}
        end
      invalid       -> invalid
    end
  end
end
