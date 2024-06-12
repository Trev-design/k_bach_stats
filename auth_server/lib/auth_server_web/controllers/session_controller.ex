defmodule AuthServerWeb.SessionController do
  use AuthServerWeb, :controller

  alias AuthServer.{Jwt, Users, Users.User}

  def refresh_token(conn, _params) do
    with %User{} = user <- Users.get_user(conn.assigns.user_id),
         {:ok, jwt, refresh}  <- Jwt.create_token_pair(user)
    do
      conn
      |> put_resp_content_type("application/json")
      |> put_resp_cookie("_auth_server_key", refresh, http_only: true, secure: true, max_age: 24*60*60, same_site: "Lax", sign: true)
      |> send_resp(200, Jason.encode!(%{name: user.name, id: user.id, jwt: jwt}))
    else
      {:error, reason} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: reason}))
      nil ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "no user found"}))
    end
  end

  def logout(conn, _params) do
    conn
    |> fetch_session()
    |> clear_session()
    |> configure_session(drop: true)
    |> delete_resp_cookie("_auth_server_key")
    |> send_resp(200, Jason.encode!(%{}))
  end
end
