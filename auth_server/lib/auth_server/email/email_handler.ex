defmodule AuthServer.Email.EmailHandler do
  alias AuthServer.Email.Mailer

  import Bamboo.Email

  use Oban.Worker, queue: :events, max_attempts: 3

  @impl Oban.Worker
  def perform(%Oban.Job{args: %{"to" => to, "name" => name}}) do
    new_email()
    |> to(to)
    |> subject("welcome #{name}")
    |> from("support@kbach.com")
    |> put_header("Reply-To", "someone@example.com")
    |> html_body("<div><p>please verify your account</p><p>your verify code is #{:crypto.strong_rand_bytes(6)}</p></div>")
    |> Mailer.deliver_now(response: false)
  end
end
