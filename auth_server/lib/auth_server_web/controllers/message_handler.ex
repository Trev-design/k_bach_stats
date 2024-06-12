defmodule AuthServerWeb.MessageHandler do
  import Plug.Conn

  require Logger

  def error_response(conn, status, message) do
    Logger.info(inspect(message))
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(status, Jason.encode!(%{message: message}))
  end

  def create_account_response(conn, credentials) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(credentials))
  end

  def signin_response(conn, credentials, token) do
    conn
    |> put_resp_content_type("application/json")
    |> put_resp_cookie("_auth_server_key", token, http_only: true, secure: true, max_age: 24*60*60, sign: true)
    |> send_resp(200, Jason.encode!(credentials))
  end

  def verify_response(conn, credentials, token) do
    conn
    |> put_resp_content_type("application/json")
    |> put_resp_cookie("_auth_server_key", token, http_only: true, secure: true, max_age: 24*60*60, sign: true)
    |> send_resp(200, Jason.encode!(credentials))
  end

  def change_password_response(conn) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(%{info: "changed password successfully"}))
  end

  def new_verify_response(conn) do
    conn
    |> put_resp_content_type("application/json")
    |> send_resp(200, Jason.encode!(%{info: "new verify successfully ordered"}))
  end
end
