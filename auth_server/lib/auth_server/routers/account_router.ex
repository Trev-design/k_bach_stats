defmodule AuthServer.Routers.AccountRouter do
  use Plug.Router

  alias AuthServer.{SessionHandler, Jwt, Schemas.Account}

  plug Plug.Logger
  plug :match
  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  plug :dispatch

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
      |> put_resp_cookie("_Refresh", refresh, http_only: true)
      |> send_resp(200, Jason.encode!(%{user:  account.user.name, jwt: jwt}))
    else
      {:error, reason} ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(500, Jason.encode!(%{message: reason}))
    end
  end
end
