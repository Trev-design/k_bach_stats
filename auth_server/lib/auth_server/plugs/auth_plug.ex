defmodule AuthServer.Plugs.AuthPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthServer.{
    Jwt,
    SessionHandler,
    Schemas.User
  }

  require Logger

  def init(opts), do: opts

  @spec call(Plug.Conn.t(), any()) :: Plug.Conn.t()
  def call(conn, _opts) do
    new_conn = fetch_cookies(conn, signed: ~w(_Refresh))
    cookie = new_conn.cookies["_Refresh"]
    if cookie == nil do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(401, Jason.encode!(%{message: "unauthorized"}))
      |> halt()
    else
      case check_auth(conn, cookie) do
        {:ok, id} -> assign(conn, :current_user_id, id)

        {:error, "token expired"} ->
          conn
          |> fetch_session()
          |> clear_session()
          |> configure_session(drop: true)
          |> delete_resp_cookie("_Refresh")
          |> send_resp(401, Jason.encode!(%{message: "session expired"}))
          |> halt()

        {:error, reason}   ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(401, Jason.encode!(%{message: reason}))
          |> halt()
      end
    end
  end

  defp check_auth(conn, cookie) do
    with {:ok, claims}     <- Jwt.check_cookie(cookie),
         {:ok, expiration} <- Map.fetch(claims, "exp"),
         {:ok, :valid}     <- check_expiration(expiration),
         {:ok, value}      <- Map.fetch(claims, "id"),
         {:ok, id}         <- check_session(conn, value)
    do
      {:ok, id}
    else
      :error -> {:error, "missing claims"}
      error  -> error
    end
  end

  defp check_expiration(expiration) do
    if expiration - Joken.current_time() <= 0 do
      {:error, "token expired"}
    else
      {:ok, :valid}
    end
  end

  defp check_session(conn, id) do
    case conn |> fetch_session() |> get_session(:session_id) do
      nil ->
        IO.inspect(conn)
        {:error, "no user validated in this session"}

      current_session ->
        check_current_user(current_session, id)
    end
  end

  defp check_current_user(current_session, id) do
    if fetch_user(current_session) == id do
      {:ok, id}
    else
      {:error, "invalid user session"}
    end
  end

  defp fetch_user(current_session) do
    with {:ok, object}          <- Jason.decode(current_session, keys: :atoms),
         {:ok, user_id}         <- Map.fetch(object, :user_id),
         {:ok , %User{} = user} <- SessionHandler.get_user(user_id)
    do
      if user.role.verified, do: user_id, else: :error
    else
      _error -> :error
    end
  end
end
