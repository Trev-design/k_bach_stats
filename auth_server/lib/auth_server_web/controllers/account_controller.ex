defmodule AuthServerWeb.AccountController do
  use AuthServerWeb, :controller

  alias AuthServerWeb.MessageHandler
  alias AuthServer.{
    Accounts,
    Accounts.Account,
    Users,
    Users.User,
    SessionHandler,
    Roles,
    Roles.Role,
    VerifyHandler,
    Jwt
  }

  def signup(conn, %{"name" => name, "email" => email, "password" => password, "confirmation" => confirmation}) do
    with true                              <- password == confirmation,
         {:ok, %Account{id: id} = account} <- Accounts.create_account(%{email: email, password_hash: password}),
         {:ok, %User{name: name} = user}   <- Users.create_user(account, %{name: name}),
         {:ok, %Role{}}                    <- Roles.create_role(user, %{verified: false, trait: "not verified"}),
         verify                            <- SessionHandler.verify_code(),
         :ok                               <- SessionHandler.send_verify_email(name, email, verify, "send_verify_email"),
         {:ok, cypher}                     <- VerifyHandler.encrypt_verification_code(id, verify),
         {:ok, "OK"}                       <- Redix.command(:verify_store, ["SET", id, cypher, "EX", 120*60])
    do
      MessageHandler.create_account_response(conn, %{id: id, name: name})

    else
      false -> MessageHandler.error_response(conn, 403, "password does not match")

      {:error, %Ecto.Changeset{errors: error_list}} -> MessageHandler.error_response(conn, 500, SessionHandler.errors(error_list))

      :error -> MessageHandler.error_response(conn, 500, "could not send email")

      {:error, %Redix.Error{message: message}} -> MessageHandler.error_response(conn, 500, message)

      {:error, %Redix.ConnectionError{reason: message}} -> MessageHandler.error_response(conn, 500, Atom.to_string(message))

      {:error, %Jason.EncodeError{message: message}} -> MessageHandler.error_response(conn, 500, message)

      {:error, reason} -> MessageHandler.error_response(conn, 500, reason)

      _invalid -> MessageHandler.error_response(conn, 500, "something went wrong")
    end
  end

  def signin(conn, %{"email" => email, "password" => password}) do
    with %Account{} = account <- Accounts.get_account_by_email(email),
         :ok                  <- SessionHandler.check_password(password, account.password_hash),
         true                 <- account.user.role.verified,
         {:ok, jwt, refresh}  <- Jwt.create_token_pair(account.user)
    do
      MessageHandler.signin_response(conn, %{name: account.user.name, jwt: jwt}, refresh)

    else
      nil -> MessageHandler.error_response(conn, 403, "invalid login credentials")

      false -> MessageHandler.error_response(conn, 403, "invalid password")

      {:error, reason} -> MessageHandler.error_response(conn, 500, reason)
    end
  end
end
