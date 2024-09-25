defmodule AuthService.Helpers do
  require Logger

  alias AuthService.Jwt

  def verify_code() do
    for _x <- 1..7 do
      :rand.uniform(9) + 48
    end
    |> List.to_string()
  end

  def errors(error_list) do
    Logger.info("your shitti ass errors are. #{error_list}")

    error_list
    |> Keyword.values()
    |> Enum.map(&elem(&1, 0))
  end

  def create_session(account, session, abo) do
    id = account.user.id
    entity = account.id
    name = account.user.name

    Jwt.create_token_pair(id, entity, name, session, abo)
  end

  def get_session_id(id) do
    case Redix.command(:user_auth_session_store, ["GET", id]) do
      {:ok, nil}            -> {:error, "invalid session"}
      {:ok, _} = session_id -> session_id
      _invalid              -> {:error, "something went wrong"}
    end
  end
end
