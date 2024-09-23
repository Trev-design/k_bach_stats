defmodule AuthServiceWeb.VerifyController do
  use AuthServiceWeb, :controller

  alias AuthServiceWeb.MessageHandler
  alias AuthService.{
    VerifyCryptoData.Access,
    Accounts,
    Accounts.Account,
    Rabbitmq,
    Helpers,
    Roles.Role,
    Roles
  }

  require Logger

  def verify(conn, %{"verify" => verify}) do
    account_id = conn.assigns[:account]

    with {:ok, cypher}         <- get_verify_cypher(account_id),
         {:ok, plain}          <- Access.decrypted(account_id, cypher),
         true                  <- verify_correct?(plain, verify),
         %Account{} = account  <- Accounts.get_full_account(account_id),
         {:ok, %Role{} = role} <- Roles.update_role(account.role, %{abo_type: "COMUNITY_USER", verified: true}),
         session_id            <- Uniq.UUID.uuid4(),
         {:ok, jwt, refresh}   <- Helpers.create_session(account, session_id, role.abo_type),
         {:ok, "OK"}           <- Redix.command(:user_auth_session_store, ["SET", account.user.id, session_id, "EX", 60 * 60 * 24]),
         {:ok, :enrolled_user} <- Rabbitmq.Access.publish_enroll_user(account, session_id, role.abo_type)
    do
      MessageHandler.session_response(conn, %{user: account.user.name, token: jwt}, refresh)

    else
      false    ->
        Logger.info("invalid verify")
      _invalid -> MessageHandler.error_response(conn, 500, %{message: "something went wrong"})
    end
  end

  def forgotten_password(conn, %{"verify" => verify, "password" => password, "confirmation" => confirmation}) do
    account_id = conn.assigns[:account]

    with {:ok, cypher}        <- get_verify_cypher(account_id),
         {:ok, plain}         <- Access.decrypted(account_id, cypher),
         true                 <- verify_correct?(plain, verify),
         %Account{} = account <- Accounts.get_full_account(account_id),
         {:ok, :verified}     <- validate_verify_status(account.role.verified),
         {:ok, %Account{}}    <- Accounts.update_account(account, %{password: password, password_confirmation: confirmation}),
         session_id           <- Uniq.UUID.uuid4(),
         {:ok, jwt, refresh}  <- Helpers.create_session(account, session_id, account.role.abo_type),
         {:ok, "OK"}          <- Redix.command(:user_auth_session_store, ["SET", account.user.id, session_id, "EX", 60 * 60 * 24]),
         {:ok, :published}    <- Rabbitmq.Access.publish_session_message(account.user.name, account.id, session_id, account.role.abo_type)
    do
      MessageHandler.session_response(conn, %{user: account.user.name, token: jwt}, refresh)

    else
      _invalid -> MessageHandler.error_response(conn, 500, "something went wrong")
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

  defp validate_verify_status(verified) do
    case verified do
      true  -> {:ok, :verified}
      false -> {:error, :not_verified}
    end
  end
end
