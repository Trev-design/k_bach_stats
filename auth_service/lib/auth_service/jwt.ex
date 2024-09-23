defmodule AuthService.Jwt do
  use Joken.Config

  require Logger

  @jwt_signer Joken.Signer.create("RS256", %{"pem" => :pem_generator.get_private_key()})
  @refresh_signer Joken.Signer.create("HS256", "a7fzdsgigj8erjre76tdsgv67fhakdhufz76ert6fdjkdhgaslalaaljdsuhfewilasdhjewuifh98efuz")

  def create_token_pair(id, entity, name, session, abo) do
    generate_token_pair(
      %{
        "id" => id,
        "entity" => entity,
        "exp" => Joken.current_time() + (24 * 60 * 60)
      },
      %{
        "id" => id,
        "entity" => entity,
        "name" => name,
        "exp" => Joken.current_time() + (60 * 15),
        "session_id" => session,
        "abo" => abo
      }
    )
  end

  def check_cookie(cookie) do
    with {:ok, claims} <- verify_and_validate(cookie, @refresh_signer),
         {:ok, expiry} <- Map.fetch(claims, "exp"),
         {:ok, :valid} <- check_cookie_expiry(expiry),
         {:ok, entity} <- Map.fetch(claims, "entity")
    do
      {:ok, entity}

    else
      :error                             -> {:error, "invalid session"}
      {:error, "cookie expired"} = error -> error
      {:error, _reason}                  -> {:error, "invalid cookie"}
    end
  end

  defp generate_token_pair(refresh_claims, jwt_claims) do
    with {:ok, refresh, _claims} <- generate_and_sign(refresh_claims, @refresh_signer),
         {:ok, jwt, _claims}     <- generate_and_sign(jwt_claims, @jwt_signer)
    do
      {:ok, jwt, refresh}

    else
      {:error, reason} ->
        Logger.info("could not create token reason: #{reason}")
        {:error, "could not create tokens"}
    end
  end

  defp check_cookie_expiry(value) do
    if value <= Joken.current_time() do
      {:error, "cookie expired"}
    else
      {:ok, :valid}
    end
  end
end
