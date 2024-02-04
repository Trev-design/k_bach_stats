defmodule AuthServer.JobHandler.EmailJob do
  def deliver(email, name) do
      %{to: email, name: name}
      |> AuthServer.Email.EmailHandler.new()
      |> Oban.insert()
  end
end
