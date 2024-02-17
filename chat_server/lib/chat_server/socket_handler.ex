defmodule ChatServer.SocketHandler do
  @behaviour :cowboy_websocket

  def init(req, _any) do
    {:cowboy_websocket, req, %{registry_key: req.path}}
  end

  def websocket_init(state) do
    Registry.register(Registry.ChatServer, state.registry_key, {})

    {:ok, state}
  end

  def websocket_handle({:text, json}, state) do
    with {:ok, _payload} <- Jason.decode(json, keys: :atoms) do
      Registry.dispatch(Registry.ChatServer ,state.registry_key, fn entries ->
        Enum.each(entries, fn {pid, _} ->
          if pid != self() do
            Process.send(pid, {:to_receive, json}, [])
          end
        end)
      end)

      {:ok, {:text, json}, state}
    end

    {:ok, {:text, Jason.encode!(%{error: "Something went wrong"})}}
  end

  def websocket_info({:to_receive, payload}, state) do
    {:reply, {:text, payload}, state}
  end
end
