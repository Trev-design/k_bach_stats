defmodule AuthServer.Jwt do
  use Joken.Config
  require Logger

  @alg "HS256"
  @secret "supersecretsinginkeywithsuperrandomlettersandnumbersandsupersecurestoredthatyoucannotstealit"
  @signer Joken.Signer.create(@alg, @secret)

  def generate_token_pair(jwt_claims, refresh_claims) do
    with {:ok, jwt, jwt_claims}         <- generate_and_sign(jwt_claims, @signer),
         {:ok, refresh, refresh_claims} <- generate_and_sign(refresh_claims, @signer)
    do
      Logger.info("jwt claims#{inspect(jwt_claims)}\nrefresh claims #{inspect(refresh_claims)}")
      {:ok, jwt, refresh}
    else
      _err -> {:error, "could not generate token"}
    end
  end
end
