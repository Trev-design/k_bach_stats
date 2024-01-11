defmodule AuthServer.Routers.SessionRouter do
  use Plug.Router

  alias AuthServer.{SessionHandler, Jwt, Schemas.User}

  plug Plug.Logger
  plug :match
  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  plug Plug.Session, store: :cookie,
    key: "_Refresh",
    signing_salt: "cookie store signing salt",
    log: :debug
  plug Plug.Session, store: :cookie,
    key: "current_user",
    signing_salt: "session store signing salt",
    log: :debug
  plug :put_secret_key_base
  plug AuthServer.Plugs.AuthPlug
  plug :dispatch

  def put_secret_key_base(conn, _) do
    put_in(conn.secret_key_base, "thekeyshouldhavemorethansixtyfourcharacterwithlettersandnumbers12345678900987654321")
  end

  get "/refresh_session" do
    IO.inspect(conn)
    with {:ok, %User{} = user}      <- SessionHandler.get_user(conn.assigns.current_user_id),
         {:ok, jwt, refresh}        <- Jwt.generate_token_pair(
                                        %{"id"   => user.id,
                                          "name" => user.name,
                                          "exp"  => Joken.current_time() + (15 * 60)},
                                        %{"id"   => user.id,
                                          "exp"  => Joken.current_time() + (24 * 60 * 60)})
    do
      conn
      |> put_resp_content_type("application/json")
      |> put_resp_cookie("_Refresh", refresh, http_only: true, secure: true, sign: true, max_age: 24*60*60)
      |> send_resp(200, Jason.encode!(%{name: user.name, id: user.id, jwt: jwt}))
    else
      {:error, reason} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: reason}))
    end
  end

  get "/signout" do
    conn
    |> fetch_session()
    |> clear_session()
    |> configure_session(drop: true)
    |> delete_resp_cookie("_Refresh")
    |> send_resp(200, Jason.encode!(%{}))
  end
end
