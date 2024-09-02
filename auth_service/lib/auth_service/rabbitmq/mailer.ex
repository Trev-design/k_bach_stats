defmodule AuthService.Rabbitmq.Mailer do
  use GenServer
  use AMQP

  @exchange "verify_email"
  @queue "verify_mail"
  @routing_key "send_verify_mail"

  def start_link(props), do: GenServer.start_link(__MODULE__, props, name: __MODULE__)

  def init(props) do
    chan =
      channel(props)
      |> setup_queue()

    {:ok, chan}
  end

  def send_verify_email(name, verify), do: GenServer.cast(__MODULE__, {:send_verify_email, name, verify})

  def handle_cast({:send_verify_email, name, verify}, channel) do
    Basic.publish(
      channel,
      @exchange,
      @routing_key,
      Jason.encode!(%{name: name, verify: verify}),
      persistent: true,
      mandatory: true
      )
  end

  defp channel(props) do
    user = Keyword.get(props, :user, "IAmTheUser")
    password = Keyword.get(props, :password, "ThisIsMyPassword")
    host = Keyword.get(props, :host, "localhost")
    port = Keyword.get(props, :port, 5672)
    vhost = Keyword.get(props, :vhost, :kbach)
    {:ok, conn} = Connection.open("amqp://#{user}:#{password}@#{host}:#{port}/#{vhost}")
    {:ok, chan} = Channel.open(conn)
    chan
  end

  defp setup_queue(chan) do
    :ok = Exchange.declare(
      chan,
      @exchange,
      :direct,
      durable: true,
      auto_delete: false,
      internal: false,
      no_wait: false
    )

    :ok = Queue.declare(
      chan,
      @queue,
      durable: true,
      auto_delete: false,
      exclusive: false,
      nowait: false)

    :ok = Queue.bind(
      chan,
      @queue,
      @exchange,
      routing_key: @routing_key,
      nowait: false
    )

    chan
  end
end
