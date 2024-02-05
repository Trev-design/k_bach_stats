defmodule AuthServer.Email.EmailHandler do
  # alias AuthServer.Email.Mailer

  import Bamboo.Email

  use Oban.Worker, queue: :events, max_attempts: 3, tags: ["email delivery"]

  @impl Oban.Worker
  def perform(%{args: %{"to" => to, "name" => name, verify: verify}}) do
    {:ok, new_email(
      to: to,
      subject: "welcome #{name}",
      from: "support@kbach.com",
      text_body: "<div><p>please verify your account</p><p>your verify code is #{verify}</p></div>"
    )}
    |> IO.inspect()
    # |> Mailer.deliver_now()
    # |> IO.inspect()
  end
end
