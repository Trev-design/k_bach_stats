defmodule AuthServer.Routers.SessionRouter do
  use Plug.Router

  plug Plug.Logger
  plug :match
  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  plug :dispatch

  get "/refresh_session" do
    conn
  end

  get "/signout" do
    conn
  end
end
