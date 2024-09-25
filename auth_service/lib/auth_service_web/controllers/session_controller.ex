defmodule AuthServiceWeb.SessionController do
  alias AuthServiceWeb.MessageHandler
  alias AuthService.Helpers

  use AuthServiceWeb, :controller

  def refresh_session(conn, _params) do
    account = conn.assigns[:account]

    with {:ok, session_id}   <- Helpers.get_session_id(account.user.id),
         {:ok, jwt, refresh} <- Helpers.create_session(account, IO.inspect(session_id), account.role.abo_type),
         {:ok, "OK"}         <- Redix.command(:user_auth_session_store, ["SET", account.user.id, session_id, "EX", 60 * 60 * 24])
    do
      IO.inspect(jwt)
      MessageHandler.session_response(conn, %{jwt: jwt}, refresh)

    else
      _invalid -> MessageHandler.error_response(conn, 403, "invalid session")
    end
  end

  def signout(conn, _params) do
    account = conn.assigns[:account]

    Redix.command(:user_auth_session_store, ["DEL", account.user.id])
    |> IO.inspect()

    conn
    |> clear_session()
    |> configure_session(drop: true)
    |> delete_resp_cookie("_auth_service_key")
    |> MessageHandler.error_response(200, "see you")
  end
end
