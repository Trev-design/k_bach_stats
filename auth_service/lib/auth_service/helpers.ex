defmodule AuthService.Helpers do
  require Logger

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
end
