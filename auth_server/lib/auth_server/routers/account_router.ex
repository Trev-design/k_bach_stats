defmodule AuthServer.Routers.AccountRouter do
  use Plug.Router

  alias AuthServer.{SessionHandler, Jwt, Schemas.Account}

  require Logger

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
  plug :dispatch

  def put_secret_key_base(conn, _) do
    put_in(conn.secret_key_base, "thekeyshouldhavemorethansixtyfourcharacterwithlettersandnumbers12345678900987654321")
  end

  post "/signin" do
    case conn.body_params do
      %{"email" => email, "password" => password} ->
        compute_signin_request(conn, email, password)
      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  post "/create" do
    Logger.info("#{inspect(conn.path_info)}")
    case conn.body_params do
      %{"name" => name, "email" => email, "password" => password, "confirmation" => confirmation} ->
        {status, response} = compute_create_request(name, email, password, confirmation)
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(status, Jason.encode!(response))

      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  defp compute_create_request(name, email, password, confirmation) do
    if password == confirmation do
      SessionHandler.register(name, email, password)
    else
      {403, %{message: "password does not match"}}
    end
  end

  defp compute_signin_request(conn, email, password) do
    with {:ok, %Account{} = account} <- SessionHandler.signin(email, password),
         {:ok, jwt, refresh}         <- Jwt.generate_token_pair(
                                          %{"id" => account.user.id,
                                            "name" => account.user.name,
                                            "exp" => Joken.current_time() + (15 * 60)},
                                          %{"id" => account.user.id,
                                            "exp" => Joken.current_time()+ (24 * 60 * 60)})
    do
      conn
      |> put_resp_content_type("application/json")
      |> fetch_session()
      |> put_session(:current_user, account.user.id)
      |> put_resp_cookie("_Refresh", refresh, http_only: true, secure: true, sign: true, max_age: 24*60*60)
      |> send_resp(200, Jason.encode!(%{name:  account.user.name, id: account.user.id, jwt: jwt}))
    else
      {:error, reason} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: reason}))
    end
  end
end
