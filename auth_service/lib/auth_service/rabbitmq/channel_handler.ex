defmodule AuthService.Rabbitmq.ChannelHandler do
  use GenServer

  require Logger

  defstruct channel: nil, pid: nil

  def start_link([]), do: GenServer.start_link(__MODULE__, [], name: __MODULE__)

  def init(init_arg), do: {:ok, init_arg}

  def set_channel(channel, pid), do: GenServer.call(__MODULE__, {:new_channel, channel, pid})
  def remove_channel(pid), do: GenServer.call(__MODULE__, {:delete_channel, pid})
  def get_all_channels(), do: GenServer.call(__MODULE__, :get_channels)

  def handle_call({:new_channel, channel, pid}, _from, state) do
    new_state = [%__MODULE__{channel: channel, pid: pid} | state]
    Logger.info("#{Enum.count(new_state)} channels so far")
    {:reply, :ok, new_state}
  end

  def handle_call({:delete_channel, pid}, _from, state) do
    object = Enum.find(state, fn object -> object.channel.pid == pid end)
    send(object.pid, {:set_to_nil, self()})
    receive do
      {:ok, :success} ->
        {:reply,
          :ok,
          Enum.reject(state, fn object ->
            object.channel.pid == pid
          end)
        }
    end
  end
  def handle_call(:get_channels, _from, state), do: {:reply, state, state}
end
