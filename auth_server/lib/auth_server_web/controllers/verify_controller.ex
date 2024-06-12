defmodule AuthServerWeb.VerifyController do
  use AuthServerWeb, :controller

  alias AuthServerWeb.MessageHandler
  alias AuthServer.{
    VerifyHandler,
    VerifyCryptoHandler,
    Roles,
    Roles.Role,
    Accounts,
    Accounts.Account,
    Jwt,
    SessionHandler
  }

  def verify(conn, %{"verification" => verification}) do
    case check_session_for_verify(conn.assigns.current_account_id, verification) do
      {:ok, credentials, refresh} ->
        MessageHandler.verify_response(conn, credentials, refresh)

      {:error, status, "session expired" = reason} ->
        VerifyCryptoHandler.purge()
        MessageHandler.error_response(conn, status, reason)

      {:error, status, reason}    ->
        MessageHandler.error_response(conn, status, reason)
    end
  end

  def change_password(conn, %{"verification" => verification, "password" => password, "confirmation" => cofirmation}) do
    case check_session_for_change_password(conn.assigns.current_account_id, verification, password, cofirmation) do
      :ok ->
        MessageHandler.change_password_response(conn)

      {:error, status, "session expired" = reason} ->
        VerifyCryptoHandler.purge()
        MessageHandler.error_response(conn, status, reason)

      {:error, status, reason}    ->
        MessageHandler.error_response(conn, status, reason)
    end
  end


  defp check_session_for_verify(account_id, verification) do
    if account_id == nil do
      {:error, 403, "invalid session id"}
    else
      with {:ok, integer}                                       <- make_integer(verification),
           {:ok, value}                                         <- Redix.command(:verify_store, ["GET", account_id]),
           {:ok, expected}                                      <- VerifyHandler.decrypt_verification_code(account_id, value),
           :ok                                                  <- check_verification(integer, expected),
           %Account{user: %{role: role, name: name}} = account  <- Accounts.get_full_account(account_id),
           {:ok, %Role{}}                                       <- Roles.change_verify_role(role, true),
           {:ok, jwt, refresh}                                  <- Jwt.create_token_pair(account.user)
      do
        {:ok, %{name: name, jwt: jwt}, refresh}
      else
        :error                                        -> {:error, 403, "invalid verification code"}
        nil                                           -> {:error, 403, "invalid session"}
        {:ok, nil}                                    -> {:error, 401, "session expired"}
        #{:error, reason}                              -> {:error, 403, reason}
        {:error, %Ecto.Changeset{errors: error_list}} -> {:error, 401, SessionHandler.errors(error_list)}
        _invalid                                      -> {:error, 500, "something went wrong"}
      end
    end
  end

  defp check_session_for_change_password(account_id, verification, password, confirmation) do
    if account_id == nil do
      {:error, 403, "invalid session id"}
    else
      with {:ok, integer}                 <- make_integer(verification),
           {:ok, value}                   <- Redix.command(:verify_store, ["GET", account_id]),
           {:ok, expected}                <- VerifyHandler.decrypt_verification_code(account_id, value),
           :ok                            <- check_verification(integer, expected),
           %Account{} = account           <- Accounts.get_account(account_id),
           :valid                         <- compare_passwords(password, confirmation),
           {:ok, %Account{}}              <- Accounts.change_password(account, Argon2.hash_pwd_salt(password))
      do
        :ok
      else
        :error                                        -> {:error, 403, "invalid verification code"}
        nil                                           -> {:error, 403, "invalid session"}
        {:ok, nil}                                    -> {:error, 401, "session expired"}
        {:error, reason}                              -> {:error, 403, reason}
        {:error, %Ecto.Changeset{errors: error_list}} -> {:error, 401, SessionHandler.errors(error_list)}
        _invalid                                      -> {:error, 500, "something went wrong"}
      end
    end
  end

  defp make_integer(verification) do
    case Integer.parse(verification) do
      {integer, ""}        -> {:ok, integer}
      {_integer, _invalid} -> {:error, "value not an integer"}
      :error               -> {:error, "invalid input"}
    end
  end

  defp check_verification(verification, excepted) do
    if verification == excepted, do: :ok, else: :error
  end

  defp compare_passwords(password, confirmation) do
    if password == confirmation, do: :valid, else: :invalid
  end
end
