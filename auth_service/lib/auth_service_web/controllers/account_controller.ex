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
         {:ok, "OK"}                                     <- Redix.command(:verify_session_store, ["SET", id, cypher, "EX", @expiry])
    do
      Logger.info("Hello #{name} here is your mega cool verify code: #{verify_code}")
      Rabbitmq.Access.publish_verify_message(verify_code, name, email, id)
      MessageHandler.create_account_response(conn, %{id: id, name: name})
    else
      {:error, _reason} -> MessageHandler.error_response(conn, 500, "Something went wrong")
      :error            -> MessageHandler.error_response(conn, 500, "Something went wrong")
    end
  end
end
