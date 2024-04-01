defmodule AuthServer.RabbitMQ do
  use GenServer

  def start_link(), do: GenServer.start_link(__MODULE__, :ok, name: __MODULE__)

  @impl GenServer
  def init(_args) do
    {:ok, conn} = AMQP.Connection.open("amqp://IAmTheUser:ThisIsMyPassword@localhost")
    {:ok, chan} = AMQP.Channel.open(conn)

    {:ok, chan}
  end

  def send_session_credentials(credentials) do
    GenServer.cast(__MODULE__, {:session_credentials, credentials})
  end

  @impl GenServer
  def handle_cast({:session_credentials, credentials}, channel) do
    {:noreply, channel}
  end
end
