defmodule AuthServiceWeb.MessageHandler do
  import Plug.Conn

  def error_response(conn, status, message) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(status, Jason.encode!(%{message: message}))
  end

  def create_account_response(conn, credentials) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(201, Jason.encode!(credentials))
  end

  def session_response(conn, credentials, refresh) do
    conn
    |> put_resp_content_type("application/json")
    |> put_resp_cookie("_auth_service_key", refresh, http_only: true, secure: true, max_age: 24*60*60, sign: true)
    |> IO.inspect()
    |> send_resp(200, Jason.encode!(credentials))
  end

  def new_verify_response(conn, message, id) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(%{success: message, id: id}))
  end
end
