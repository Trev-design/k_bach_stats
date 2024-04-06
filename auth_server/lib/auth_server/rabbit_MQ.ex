defmodule AuthServer.RabbitMQ do
  use GenServer

  def start_link() do
    GenServer.start_link(__MODULE__, [], name: __MODULE__)
  end

  @impl GenServer
  def init(_opts) do
    {:ok, conn} = AMQP.Connection.open([])
    {:ok, channel} = AMQP.Channel.open(conn)

    {:ok, channel}
  end

  def send_new_user_credentials(credentials) do
    GenServer.cast(__MODULE__, {:new_user, credentials})
  end

  @impl GenServer
  def handle_cast({:new_user, credentials}, channel) do
    AMQP.Basic.publish(channel, "session_start", "start_session.register", credentials)
    {:noreply, channel}
  end
end
