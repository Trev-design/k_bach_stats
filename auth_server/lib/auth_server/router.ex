defmodule AuthServer.Router do
  use Plug.Router

  plug Plug.Logger
  plug :match
  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  plug :dispatch

  get "/" do
    send_resp(conn, 200, "OK")
  end

  forward "/account", to: AuthServer.Routers.AccountRouter
  forward "/session", to: AuthServer.Routers.SessionRouter

  if Mix.env() == :dev do
    forward "/sent_emails", to: Bamboo.SentEmailViewerPlug
  end

  match _ do
    send_resp(conn, 500, "ERROR")
  end
end
