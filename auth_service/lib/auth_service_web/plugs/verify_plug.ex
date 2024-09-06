defmodule AuthServiceWeb.VerifyPlug do
  @behaviour Plug

  import Plug.Conn

  alias AuthService.{Accounts, Accounts.Account}

  def init(opts), do: opts

  def call(conn, _) do
    with [uuid]     <- get_req_header(conn, "userid"),
         %Account{} <- Accounts.get_account(uuid)
    do
      assign(conn, :account, uuid)

    else
      _invalid ->
        conn
        |> put_resp_content_type("application/json")
        |> send_resp(401, Jason.encode!(%{message: "invalid session"}))
        |> halt()
    end
  end
end
