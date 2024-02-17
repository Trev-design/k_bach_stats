defmodule ChatServer.SocketHandler do
  @behaviour :cowboy_websocket

  def init(req, _any) do
    {:cowboy_websocket, req, %{registry_key: req.path}}
  end

  def websocket_init(state) do
    Registry.register(Registry.ChatServer, state.registry_key, {})

    {:ok, state}
  end

  def websocket_handle(frame = {:text, json}, state) do
    with {:ok, _payload} <- Jason.decode(json, keys: :atoms) do
      distribute_messages(self(), json, state.registry_key)
      {[frame], state}
    end

    {[{:text, Jason.encode!(%{error: "Something went wrong"})}], state}
  end

  def websocket_info({:to_distribute, payload}, state) do
    {[{:text, payload}], state}
  end

  defp distribute_messages(sender, message, registry_key) do
    Task.Supervisor.start_child(ChatServer.TaskSupervisor, fn ->
      Registry.dispatch(Registry.ChatServer, registry_key, fn entries ->
        Enum.each(entries, fn {pid, _} ->
          if pid != sender do
            Process.send(pid, {:to_distribute, message}, [])
          end
        end)
      end)
    end)
  end
end
