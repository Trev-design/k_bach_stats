defmodule AuthService.VerifyCryptoData.PurgeHelper do
  require Logger

  alias AuthService.VerifyCryptoData.HandlerFunctions
  def prepare_purge() do
    Task.Supervisor.start_child(PurgeHelper.Supervisor, fn -> make_purge() end)
  end

  defp make_purge() do
    case HandlerFunctions.get_deletables_transaction() do
      nil   ->
        Logger.info("No messages to delete")
        {:ok, :no_deletables}
      terms ->
        Logger.info("yaaaaaaahoooooooooo")
        terms
        |> Stream.chunk_every(10)
        |> Stream.map(fn chunk -> handle_purge(chunk) end)
        |> Stream.each(fn task -> Task.await(task) end)
        |> Stream.run()
    end
  end

  defp handle_purge(terms) do
    Task.async(fn ->
      Poolex.run(
        :purge,
        fn pid ->
          GenServer.call(pid, {:purge, terms})
        end
      )
    end)
  end
end
