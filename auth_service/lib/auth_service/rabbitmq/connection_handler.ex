defmodule AuthService.Rabbitmq.ConnectionHandler do
  alias AuthService.Rabbitmq.ChannelHandler
  alias AuthService.Rabbitmq.HandlerFunctions
  use GenServer

  def start_link(props), do: GenServer.start_link(__MODULE__, props, name: :rmq_connection)

  def init(props) do
    name = Keyword.get(props, :name, "IAmTheUser")
    password = Keyword.get(props, :password, "ThisIsMyPassword")
    host = Keyword.get(props, :host, "localhost")
    port = Keyword.get(props, :port, 5672)
    vhost = Keyword.get(props, :vhost, "kbach")
    {:ok, connection} = HandlerFunctions.setup_connections("amqp://#{name}:#{password}@#{host}:#{port}/#{vhost}")

    {:ok, connection}
  end

  def handle_info({:declare_channel, pid}, connection) do
    send(pid, {:get_connection, connection})
    {:noreply, connection}
  end

  def handle_info({:stop_channel, channel}, connection) do
    ChannelHandler.remove_channel(channel.pid)
    HandlerFunctions.close_channel(channel)
    {:noreply, connection}
  end

  def terminate(reason, connection) do
    ChannelHandler.get_all_channels()
    |> Stream.each(fn channel -> HandlerFunctions.close_channel(channel) end)
    |> Stream.run()

    HandlerFunctions.close_connection(connection)

    {:stop, reason}
  end
end
