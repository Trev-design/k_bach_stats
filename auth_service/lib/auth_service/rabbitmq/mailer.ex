defmodule AuthService.Rabbitmq.Mailer do
  use GenServer

  alias AuthService.Rabbitmq.HandlerFunctions

  @exchange "verify_email"
  @queue "verify_mail"
  @routing_key "send_verify_mail"

  def start_link(props), do: GenServer.start_link(__MODULE__, props, name: __MODULE__)

  def init(props) do
    user = Keyword.get(props, :user, "IAmTheUser")
    password = Keyword.get(props, :password, "ThisIsMyPassword")
    host = Keyword.get(props, :host, "localhost")
    port = Keyword.get(props, :port, 5672)
    vhost = Keyword.get(props, :vhost, :kbach)

    chan =
      HandlerFunctions.setup_connections("amqp://#{user}:#{password}@#{host}:#{port}/#{vhost}")
      |> HandlerFunctions.declare_exchange(@exchange, true)
      |> HandlerFunctions.declare_queue(@queue, true, false)
      |> HandlerFunctions.bind_queue(@exchange, @routing_key, @queue)

    {:ok, chan}
  end

  def send_verify_email(name, verify), do: GenServer.cast(__MODULE__, {:send_verify_email, name, verify})

  def handle_cast({:send_verify_email, name, verify}, channel) do
    HandlerFunctions.publish(channel, @exchange, @routing_key, Jason.encode!(%{name: name, verify: verify}))

    {:noreply, channel}
  end
end
