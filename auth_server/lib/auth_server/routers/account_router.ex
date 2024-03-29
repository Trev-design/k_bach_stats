defmodule AuthServer.Routers.AccountRouter do

  use Plug.Router

  alias AuthServer.{
    SessionHandler,
    Schemas.Account,
    Schemas.User,
    Routers.RouterHelpers
  }

  plug Corsica,
    origins: "*",
    allow_headers: ["accept", "content-type", "authorization"],
    allow_methods: ["GET", "POST", "OPTIONS"],
    allow_credentials: true

  plug Plug.Logger

  plug :match

  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason

  plug Plug.Session, store: :cookie,
    key: "_verify",
    signing_salt: "averycomplicatedandstrongsigningsaltwithalotofcharactersandnumbers1234567890",
    log: :debug

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
        compute_new_verify_request(conn, email, "send_verify_email")

      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  post "/forgotten_password" do
    case conn.body_params do
      %{"email" => email} ->
        compute_new_verify_request(conn, email, "send_forgot_password_email")

      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
    end
  end

  put "/change_password/:id" do
    case conn.body_params do
      %{"verify" => verify, "password" => password, "confirmation" => confirmation} ->
        compute_change_password_request(conn, verify, password, confirmation)

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
        verification_code_response(conn, id, name, email, "send_verify_email")

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
    new_conn = fetch_cookies(conn, signed: ~w(_verify))
    cookie = new_conn.cookies["_verify"]
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
        |> put_session(:session_id, Jason.encode!(%{user_id: user.id, session: UUID.uuid4()}))
        |> put_resp_cookie("_Refresh", refresh, http_only: true, secure: true, sign: true, max_age: 24*60*60, same_site: "Strict")
        |> delete_resp_cookie("_verify")
        |> send_resp(200, Jason.encode!(%{name:  user.name, id: user.id, jwt: jwt}))
      else
        {:error, reason} ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(500, Jason.encode!(%{message: reason}))
      end
    end
  end

  defp compute_new_verify_request(conn, email, route) do
    case SessionHandler.get_by_email(email) do
      %Account{id: account_id, email: email, user: %{name: name, id: user_id}} ->
        cond do
          route == "send_forgot_password_email" ->
            verification_code_response(conn, account_id, name, email, route)

          route == "send_verify_email" ->
            verification_code_response(conn, user_id, name, email, route)

          true ->
            conn
            |> put_resp_content_type("application/json")
            |> send_resp(500, Jason.encode!(%{message: "something went wrong"}))
        end
      nil ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(403, Jason.encode!(%{message: "invalid email"}))
    end
  end

  defp compute_change_password_request(conn, verify, password, confirmation) do
    new_conn = fetch_cookies(conn, signed: ~w(_verify))
    cookie = new_conn.cookies["_verify"]
    if cookie == nil do
      conn
      |> put_resp_content_type("application/json")
      |> send_resp(401, Jason.encode!(%{message: "unauthorized"}))
    else
      with {:ok, %Account{} = account} <- SessionHandler.check_email_verification(verify, cookie),
           true                        <- password == confirmation,
           {:ok, %Account{}}           <- SessionHandler.update_password(account, password)
      do
        conn
        |> put_resp_content_type("application/json")
        |> fetch_session()
        |> delete_resp_cookie("_verify")
        |> send_resp(200, Jason.encode!(%{message: "changed password successful"}))
      else
        {:error, error} ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(500, Jason.encode!(%{message: error}))

        false ->
          conn
          |> put_resp_content_type("application/json")
          |> send_resp(500, Jason.encode!(%{message: "password does not match"}))
      end
    end
  end

  defp verification_code_response(conn, id, name, email, route) do
    random_number = RouterHelpers.create_verify_code()

    mailer_request(email, "hello #{name} your verify code is #{random_number}", route)

    conn
    |> put_resp_content_type("application/json")
    |> put_resp_cookie(
      "_verify",
      Jason.encode!(%{id: id, name: name, verify: random_number}),
      http_only: true,
      secure: true,
      sign: true,
      max_age: 60*60,
      same_site: "Strict"
    )
    |> send_resp(200, Jason.encode!(%{guest: name}))
  end

  defp mailer_request(email, data, route) do
    Task.Supervisor.start_child(
      MailRequest.Supervisor,
      fn ->
        IO.puts("try to send email")
        HTTPoison.start()
        with {:ok, %HTTPoison.AsyncResponse{}} <- make_email_response(email, data, route) do
          receive do
            %HTTPoison.AsyncStatus{code: 200} ->
              IO.puts("received an ok")
              :ok

            %HTTPoison.AsyncStatus{code: _invalid} ->
              IO.puts("received an error")
              :error
          end
          IO.puts("finished")
        else
          _err -> :error
        end
      end)
  end

  defp make_email_response(email, data, route) do
    HTTPoison.post(
      "127.0.0.1:8080/#{route}",
      Jason.encode!(%{
        to: email,
        subject: "your verify code",
        data: data}),
      ["Content-Type": "application/json"],
      [stream_to: self(), recv_timeout: 5 * 60 * 1000])
  end
end
