defmodule ChatServer.Router do
  use Plug.Router

  plug :match
  plug Plug.Parsers,
    parsers: [:json],
    pass: ["application/json"],
    json_decoder: Jason
  plug :dispatch

  get "/" do
    send_resp(conn, 200, Jason.encode!(%{hello: "world"}))
  end

  match _ do
    send_resp(conn, 404, Jason.encode!(message: "not found"))
  end
end
