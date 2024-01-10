defmodule AuthServer.Routers.SessionRouter do
  use Plug.Router

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
    send_resp(conn, 200, Jason.encode!(%{popo: "pippi"}))
  end

  get "/signout" do
    send_resp(conn, 200, Jason.encode!(%{popo: "pippi"}))
  end
end
