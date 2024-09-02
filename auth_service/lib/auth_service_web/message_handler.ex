defmodule AuthServiceWeb.MessageHandler do
  import Plug.Conn

  def error_response(conn, status, message) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(status, Jason.encode!(%{message: message}))
  end

  def create_account_response(conn, account, username) do
    conn
    |> put_resp_content_type("application/json")
    |> put_resp_header("account", account)
    |> send_resp(201, Jason.encode!(%{user: username}))
  end
end
