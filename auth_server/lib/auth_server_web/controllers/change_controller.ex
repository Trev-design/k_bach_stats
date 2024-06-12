defmodule AuthServerWeb.ChangeController do
  use AuthServerWeb, :controller

  alias AuthServerWeb.MessageHandler
  alias AuthServer.{
    SessionHandler,
    Accounts,
    Accounts.Account,
    Users.User
}

  def new_verify(conn, %{"email" => email}) do
    with %Account{user: %User{name: name}} <- Accounts.get_account_by_email(email),
         :ok                               <- SessionHandler.send_verify_email(name, email, SessionHandler.verify_code(), "send_verify_email")
    do
      MessageHandler.new_verify_response(conn)
    else
      nil    -> MessageHandler.error_response(conn, 500, "could not find account")
      :error -> MessageHandler.error_response(conn, 500, "could not sent verify code")
    end
  end

  def forgotten_password(conn, %{"email" => email}) do
    with %Account{user: %User{name: name}} <- Accounts.get_account_by_email(email),
         :ok                               <- SessionHandler.send_verify_email(name, email, SessionHandler.verify_code(), "send_verify_email")
    do
      MessageHandler.new_verify_response(conn)
    else
      nil    -> MessageHandler.error_response(conn, 500, "could not find account")
      :error -> MessageHandler.error_response(conn, 500, "could not sent verify code")
    end
  end
end
