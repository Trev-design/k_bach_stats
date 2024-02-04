defmodule AuthServer.JobHandler.EmailJob do
  def deliver(email, name) do
    job =
      %{"email" => email, "name" => name}
      |> AuthServer.Email.EmailHandler.new()
      |> Oban.insert()

    IO.inspect(job)
  end
end
