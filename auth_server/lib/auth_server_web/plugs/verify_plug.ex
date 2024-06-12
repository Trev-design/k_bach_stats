defmodule AuthServerWeb.VerifyPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthServer.{Accounts, Accounts.Account}

  require Logger

  def init(opts), do: opts

  @spec call(Plug.Conn.t(), any()) :: Plug.Conn.t()
  def call(conn, _opts) do
    Logger.info(inspect(conn))

    with [uuid]     <- get_req_header(conn, "userid"),
         %Account{} <- Accounts.get_account(uuid)
    do
      assign(conn, :current_account_id, uuid)

    else
      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(401, Jason.encode!(%{message: "invalid session"}))
        |> halt()
    end
  end
end
