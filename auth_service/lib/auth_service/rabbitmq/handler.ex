defmodule AuthService.Rabbitmq.Handler do
  use GenServer

  alias AuthService.Rabbitmq.HandlerFunctions

  @exchange "verify"
  @exchange2 "session"
  @queue "verification_email"
  @queue2 "user_session"
  @routing_key1 "send_verify_email"
  @routing_key2 "senf_forgotten_password_email"
  @routing_key3 "session_distribute"

  def start_link(props), do: GenServer.start_link(__MODULE__, props, name: __MODULE__)

  def init(props) do
    user = Keyword.get(props, :user, "IAmTheUser")
    password = Keyword.get(props, :password, "ThisIsMyPassword")
    host = Keyword.get(props, :host, "localhost")
    port = Keyword.get(props, :port, 5672)
    vhost = Keyword.get(props, :vhost, "kbach")

    chan =
      HandlerFunctions.setup_connections("amqp://#{user}:#{password}@#{host}:#{port}/#{vhost}")
      |> HandlerFunctions.declare_exchange(@exchange, true)
      |> HandlerFunctions.declare_exchange(@exchange2, true)
      |> HandlerFunctions.declare_queue(@queue, true, false)
      |> HandlerFunctions.declare_queue(@queue2, true, false)
      |> HandlerFunctions.bind_queue(@exchange, @routing_key1, @queue)
      |> HandlerFunctions.bind_queue(@exchange, @routing_key2, @queue)
      |> HandlerFunctions.bind_queue(@exchange2, @routing_key3, @queue2)

    {:ok, chan}
  end

  def send_verify_email(user_id, email, name, verify) do
    GenServer.cast(__MODULE__, {:send_verify_email, user_id, email, name, verify})
  end

  def create_session(user_id, name, session_id) do
    GenServer.cast(__MODULE__, {:create_session, user_id, name, session_id})
  end

  def handle_cast({:send_verify_email, user_id, email, name, verify}, channel) do
    HandlerFunctions.publish(
      channel,
      @exchange,
      @routing_key1,
      Jason.encode!(%{user_id: user_id, email: email, name: name, verify: verify, kind: "verify"})
    )

    {:noreply, channel}
  end

  def handle_cast({:create_session, user_id, name, session_id}, channel) do
    HandlerFunctions.publish(
      channel,
      @exchange2,
      @routing_key3,
      Jason.encode!(%{user_id: user_id, name: name, session_id: session_id})
    )
  end
end
