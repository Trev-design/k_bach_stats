defmodule AuthServiceWeb.AccountController do
  use AuthServiceWeb, :controller

  require Logger

  alias AuthServiceWeb.MessageHandler
  alias AuthService.{
    Accounts,
    Accounts.Account,
    Users,
    Users.User,
    Roles,
    Roles.Role,
    Helpers,
    VerifyCryptoData,
    Rabbitmq
  }

  @expiry 60 * 120

  def signup(conn, %{"name" => name, "email" => email, "password" => password, "confirm" => confirm}) do
    with {:ok, %Account{id: id, email: email} = account} <- Accounts.create_account(%{email: email, password: password, password_confirmation: confirm}),
         {:ok, %User{name: name}}                        <- Users.create_user(account, %{name: name}),
         {:ok, %Role{}}                                  <- Roles.create_role(account, %{}),
         verify_code                                     <- Helpers.verify_code(),
         {:ok, cypher}                                   <- VerifyCryptoData.Access.encrypted(id, Jason.encode!(%{id: id, verify: verify_code})),
         {:ok, "OK"}                                     <- Redix.command(:verify_session_store, ["SET", id, cypher, "EX", @expiry]),
         {:ok, :published}                               <- Rabbitmq.Access.publish_verify_message(verify_code, name, email, id)
    do
      Logger.info("Hello #{name} here is your mega cool verify code: #{verify_code}")
      MessageHandler.create_account_response(conn, %{id: id, name: name})
    else
      {:error, _reason} -> MessageHandler.error_response(conn, 500, "Something went wrong")
      :error            -> MessageHandler.error_response(conn, 500, "Something went wrong")
    end
  end

  def signin(conn, %{"email" => email, "password" => password}) do
    with %Account{} = account <- Accounts.get_full_account_by_email(email),
         true                 <- account.role.verified,
         true                 <- Argon2.verify_pass(password, account.password_hash),
         session_id           <- Uniq.UUID.uuid4(),
         {:ok, jwt, refresh}  <- Helpers.create_session(account, session_id, account.role.abo_type),
         {:ok, :published}    <- Rabbitmq.Access.publish_session_message(account.user.name, account.id, session_id, account.role.abo_type)
    do
      MessageHandler.session_response(conn, %{user: account.user.name, jwt: jwt}, refresh)

    else
      nil   -> MessageHandler.error_response(conn, 401, "invalid login credentials")
      false -> MessageHandler.error_response(conn, 401, "invalid login credentials")
      _else -> MessageHandler.error_response(conn, 500, "something, went wrong")
    end
  end

  def request_new_verify(conn, %{"email" => email, "kind" => kind}) do
    with %Account{} = account <- Accounts.get_full_account_by_email(email),
         verify_code          <- Helpers.verify_code(),
         {:ok, cypher}        <- VerifyCryptoData.Access.encrypted(account.id, verify_code),
         {:ok, "OK"}          <- Redix.command(:verify_session_store, ["SET", account.id, cypher, "EX", @expiry]),
         {:ok, :published}    <- make_delivery(kind, verify_code, account.user.name, email, account.id)
    do
      MessageHandler.new_verify_response(conn, "new verify code is sent to your email", account.id)

    else
      nil -> MessageHandler.error_response(conn, 401, "invalid credentials")
      _   -> MessageHandler.error_response(conn, 500, "something went wrong")
    end
  end

  defp make_delivery(kind, verify_code, username, email, account_id) do
    case kind do
      "new_verify" ->
        Rabbitmq.Access.publish_verify_message(verify_code, username, email, account_id)

      "forgotten_password" ->
        Rabbitmq.Access.publish_forgotten_password_message(verify_code, username, email, account_id)
    end
  end
end
