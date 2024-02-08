defmodule AuthServer.Routers.AccountRouter do
  use Plug.Router

  alias AuthServer.{
    SessionHandler,
    Schemas.Account,
    Schemas.User,
    JobHandler.EmailJob,
    Routers.RouterHelpers
  }

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
    key: "session_id",
    signing_salt: "session store signing salt",
    http_only: true,
    secure: true,
    sign: true,
    same_site: "Strict",
    log: :debug

  plug Plug.Session, store: :ets,
    table: :verify_session_table,
    key: "_verify",
    signing_salt: "averycomplicatedandstrongsigningsaltwithalotofcharactersandnumbers1234567890",
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
    case conn.body_params do
      %{"name" => name, "email" => email, "password" => password, "confirmation" => confirmation} ->
        make_response(conn, name, email, password, confirmation)

      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  post "/verify" do
    case conn.body_params do
      %{"verification" => verification} ->
        compute_verification(conn, verification)

      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  post "/new_verify" do
    case conn.body_params do
      %{"email" => email} ->
        compute_new_verify_request(conn, email)

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

  defp make_response(conn, name, email, password, confirmation) do
    case compute_create_request(name, email, password, confirmation) do
      {200, %{id: id}} ->
        verification_code_response(conn, id, name, email)

      {status, response} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(status, Jason.encode!(response))
    end
  end

  defp compute_signin_request(conn, email, password) do
    with {:ok, %Account{} = account} <- SessionHandler.signin(email, password),
         {:ok, jwt, refresh}         <- RouterHelpers.add_claims_and_generate(account.user.id, account.user.name)
    do
      conn
      |> put_resp_content_type("application/json")
      |> fetch_session()
      |> put_session(:session_id, Jason.encode!(%{user_id: account.user.id, session: UUID.uuid4()}))
      |> put_resp_cookie("_Refresh", refresh, http_only: true, secure: true, sign: true, max_age: 24*60*60, same_site: "Strict")
      |> send_resp(200, Jason.encode!(%{name:  account.user.name, id: account.user.id, jwt: jwt}))
    else
      {:error, reason} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: reason}))
    end
  end

  defp compute_verification(conn, verification) do
    new_conn = fetch_cookies(conn, signed: ~w(_verification))
    cookie = new_conn.cookies["_verification"]
    if cookie == nil do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(401, Jason.encode!(%{message: "unauthorized"}))
    else
      with {:ok, %User{} = user} <- SessionHandler.check_verification(verification, cookie),
           {:ok, jwt, refresh}   <- RouterHelpers.add_claims_and_generate(user.id, user.name)
      do
        conn
        |> put_resp_content_type("application/json")
        |> fetch_session()
        |> delete_resp_cookie("_verify")
        |> put_session(:session_id, Jason.encode!(%{user_id: user.id, session: UUID.uuid4()}))
        |> put_resp_cookie("_Refresh", refresh, http_only: true, secure: true, sign: true, max_age: 24*60*60, same_site: "Strict")
        |> send_resp(200, Jason.encode!(%{name:  user.name, id: user.id, jwt: jwt}))
      else
        {:error, reason} ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(500, Jason.encode!(%{message: reason}))
      end
    end
  end

  defp compute_new_verify_request(conn, email) do
    case SessionHandler.get_by_email(email) do
      %Account{email: email, user: %{name: name, id: id}} ->
        verification_code_response(conn, id, name, email)

      nil ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(403, Jason.encode!(%{message: "invalid email"}))
    end
  end

  defp verification_code_response(conn, id, name, email) do
    random_number =
      for _x <- 1..7 do
        :rand.uniform(9) + 48
      end
      |> List.to_integer()

      EmailJob.deliver(email, name, random_number)
      |> IO.inspect()

      conn
      |> put_resp_content_type("application/json")
      |> put_resp_cookie(
        "_verify",
        Jason.encode(%{id: id, name: name, verify: random_number}),
        http_only: true,
        secure: true,
        sign: true,
        max_age: 60*60,
        same_site: "Strict"
      )
      |> send_resp(200, Jason.encode!(%{guest: name}))
  end
end
