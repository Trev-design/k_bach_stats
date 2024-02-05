defmodule AuthServer.JobHandler.EmailJob do
  def deliver(email, name, verify) do
      %{to: email, name: name, verify: verify}
      |> AuthServer.Email.EmailHandler.new()
      |> Oban.insert()
  end
end
