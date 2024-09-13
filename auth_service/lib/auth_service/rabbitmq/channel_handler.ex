defmodule AuthService.Rabbitmq.ChannelHandler do
  use GenServer

  def start_link([]), do: GenServer.start_link(__MODULE__, [], name: __MODULE__)

  def init(init_arg), do: {:ok, init_arg}

  def set_channel(channel), do: GenServer.call(__MODULE__, {:new_channel, channel})
  def remove_channel(pid), do: GenServer.call(__MODULE__, {:delete_channel, pid})
  def get_all_channels(), do: GenServer.call(__MODULE__, :get_channels)

  def handle_call({:new_channel, channel}, _from, state), do: {:reply, :ok, [channel | state]}
  def handle_call({:delete_channel, pid}, _from, state) do
    {
      :reply,
      :ok,
      Enum.reject(state, fn %AMQP.Channel{pid: channel_pid} -> channel_pid == pid end)
    }
  end
  def handle_call(:get_channels, _from, state), do: {:reply, state, state}
end
