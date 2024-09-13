defmodule AuthServiceWeb.VerifyController do
  use AuthServiceWeb, :controller

  alias AuthServiceWeb.MessageHandler
  alias AuthService.{
    VerifyCryptoData.Access,
    Accounts,
    Accounts.Account,
    Jwt
  }

  require Logger

  def verify(conn, %{"verify" => verify}) do
    account_id = conn.assigns[:account]

    with {:ok, cypher}        <- get_verify_cypher(account_id),
         {:ok, plain}         <- Access.decrypted(account_id, cypher),
         true                 <- verify_correct?(plain, verify),
         %Account{} = account <- Accounts.get_full_account(account_id),
         {:ok, jwt, refresh}  <- create_session(account)
    do
      MessageHandler.session_response(conn, %{user: account.user.name, token: jwt}, refresh)

    else
      false    ->
        Logger.info("invalid verify")
      _invalid -> MessageHandler.error_response(conn, 500, %{message: "something went wrong"})
    end
  end

  defp create_session(account) do
    id = account.user.id
    name = account.user.name
    session = Uniq.UUID.uuid4()

    case Jwt.create_token_pair(id, name, session) do
      {:ok, _jwt, _refresh} = result->
        result

      invalid    -> invalid
    end
  end

  defp verify_correct?(plain, verify), do: plain == verify

  defp get_verify_cypher(user_id) do
    case Redix.command(:verify_session_store, ["GET", user_id]) do
      {:ok, nil}      -> {:error, :expired}
      {:ok, _} = item -> item
      _invalid        -> {:error, :something_went_wrong}
    end
  end
end
