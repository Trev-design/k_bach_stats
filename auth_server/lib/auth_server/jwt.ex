defmodule AuthServer.Jwt do
  use Joken.Config

  @alg "HS256"
  @secret "supersecretsinginkeywithsuperrandomlettersandnumbersandsupersecurestoredthatyoucannotstealit"
  @signer Joken.Signer.create(@alg, @secret)

  def generate_token_pair(jwt_claims, refresh_claims) do
    with {:ok, jwt, _jwt_claims}         <- generate_and_sign(jwt_claims, @signer),
         {:ok, refresh, _refresh_claims} <- generate_and_sign(refresh_claims, @signer)
    do
      {:ok, jwt, refresh}
    else
      _err -> {:error, "could not generate token"}
    end
  end

  def generate_expired_cookie(), do: generate_and_sign(%{"exp" => 0}, @signer)

  def check_cookie(cookie) do
    with {:ok, claims} <- verify_and_validate(cookie, @signer),
         {:ok, value}  <- Map.fetch(claims, "exp")
    do
      check_cookie_expiry(claims, value)
    else
      _invalid -> {:error, "invalid cookie"}
    end
  end

  defp check_cookie_expiry(claims, value) do
    if value <= Joken.current_time() do
      {:error, "expired cookie"}
    else
      {:ok, claims}
    end
  end
end
