defmodule AuthServerWeb.AuthPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthServer.{Jwt}

  require Logger

  def init(opts), do: opts

  @spec call(Plug.Conn.t(), any()) :: Plug.Conn.t()
  def call(conn, _opts) do
    new_conn = fetch_cookies(conn, signed: ~w(_auth_server_key))
    cookie = new_conn.cookies["_auth_server_key"]
    Logger.info("the cookie is: #{cookie}")
    if cookie == nil do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(401, Jason.encode!(%{message: "unauthorized"}))
      |> halt()

    else
      case check_auth(cookie) do
        {:ok, user_id} -> assign(conn, :user_id, user_id)

        {:error, "token expired"} ->
          conn
          |> clear_session()
          |> configure_session(drop: true)
          |> delete_resp_cookie("_Refresh")
          |> put_resp_content_type("application/json")
          |> send_resp(401, Jason.encode!(%{message: "session expired"}))
          |> halt()

        {:error, reason} ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(401, Jason.encode!(%{message: reason}))
          |> halt()
      end
    end
  end

  defp check_auth(cookie) do
    with {:ok, claims}           <- Jwt.check_cookie(cookie),
         {:ok, expiration}       <- Map.fetch(claims, "exp"),
         :ok                     <- check_expiration(expiration),
         {:ok, _user_id} = valid <- Map.fetch(claims, "id")
    do
      valid
    else
      :error -> {:error, "invalid credentials"}
      error  -> error
    end
  end

  defp check_expiration(expiration) do
    if expiration - Joken.current_time() <= 0 do
      {:error, "token expired"}
    else
      :ok
    end
  end
end
