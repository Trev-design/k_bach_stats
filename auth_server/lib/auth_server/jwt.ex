defmodule AuthServer.Jwt do
  use Joken.Config

  require Logger

  @alg "HS256"
  @secret "123456789abcdef0fedcba987654321halliillahhalloollah123456789abcdef0fedcba987654321"
  @signer Joken.Signer.create(@alg, @secret)

  def create_token_pair(user) do
    generate_token_pair(
      %{
        "id" => user.id,
        "exp" => Joken.current_time() + (24 * 60 * 60)
      },
      %{
        "id" => user.id,
        "name" => user.name,
        "exp" => Joken.current_time() + (15 * 60)
      })
  end

  defp generate_token_pair(refresh_claims, token_claims) do
    with {:ok, jwt, _claims}     <- generate_and_sign(token_claims, @signer),
         {:ok, refresh, _claims} <- generate_and_sign(refresh_claims, @signer)
    do
      {:ok, jwt, refresh}
    else
      {:error, _reason} -> {:error, "could not create tokens"}
    end
  end

  def check_cookie(cookie) do
    with {:ok, claims} <- verify_and_validate(cookie, @signer),
         {:ok, expiry} <- Map.fetch(claims, "exp")
    do
      Logger.info("the claims are #{inspect(claims)}")
      check_cookie_expiry(claims, expiry)
    else
      _invalid -> {:error, "invalid cookie"}
    end
  end

  defp check_cookie_expiry(claims, value) do
    if value <= Joken.current_time() do
      {:error, "cookie expired"}
    else
      {:ok, claims}
    end
  end
end
