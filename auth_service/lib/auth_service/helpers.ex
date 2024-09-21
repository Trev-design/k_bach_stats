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
    name = account.user.name

    Jwt.create_token_pair(id, name, session, abo)
  end
end
