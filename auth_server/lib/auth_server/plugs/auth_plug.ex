defmodule AuthServer.Plugs.AuthPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthServer.Jwt

  require Logger

  def init(opts), do: opts

  def call(conn, _opts) do
    new_conn = fetch_cookies(conn, signed: ~w(_Refresh))
    cookie = new_conn.cookies["_Refresh"]
    IO.inspect(new_conn)
    if cookie == nil do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(403, Jason.encode!(%{message: "unauthorized"}))
      |> halt()
    else
      case check_auth(conn, cookie) do
        {:ok, :authorized} -> conn

        {:error, reason}   ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(403, Jason.encode!(%{message: reason}))
          |> halt()
      end
    end
  end

  defp check_auth(conn, cookie) do
    with {:ok, claims}      <- Jwt.check_cookie(cookie),
         {:ok, value}       <- Map.fetch(claims, "id"),
         {:ok, :valid_user} <- check_session(conn, value)
    do
      {:ok, :authorized}
    else
      :error -> {:error, "missing claims"}
      error  -> error
    end
  end

  defp check_session(conn, id) do
    case conn |> fetch_session() |> get_session(:current_user) do
      nil ->
        {:error, "no user validated in this session"}

      current_session ->
        check_current_user(current_session, id)
    end
  end

  defp check_current_user(current_session, id) do
    if current_session == id, do: {:ok, :valid_user}, else: {:error, "invalid user session"}
  end
end
