defmodule AuthServiceWeb.AuthPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthService.Jwt
  alias AuthServiceWeb.MessageHandler
  alias AuthService.{Accounts, Accounts.Account}

  def init(opts), do: opts

  def call(conn, _opts) do
    new_conn = fetch_cookies(conn, signed: ~w(_auth_service_key))
    cookie = new_conn.cookies["_auth_service_key"]

    if cookie == nil do
      MessageHandler.error_response(conn, 401, "invalid session") |> halt()
    else
      case check_auth(cookie) do
        {:ok, account} -> assign(conn, :account, account)

        {:error, "cookie expired"} ->
          conn
          |> clear_session()
          |> configure_session(drop: true)
          |> delete_resp_cookie("_auth_service_key")
          |> MessageHandler.error_response(401, "cookie expired")
          |> halt()

        {:error, reason} ->
          MessageHandler.error_response(conn, 401, reason) |> halt()
      end
    end
  end

  defp check_auth(cookie) do
    with {:ok, entity}        <- Jwt.check_cookie(cookie),
         %Account{} = account <- Accounts.get_full_account(entity),
         true                 <- account.role.verified
    do
      {:ok, account}

    else
      nil     -> {:error, "invalid session"}
      false   -> {:error, "invalid session"}
      invalid -> invalid
    end
  end
end
