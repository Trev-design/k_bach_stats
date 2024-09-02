defmodule AuthService.VerifyCryptoData.PurgeHelper do

  alias AuthService.VerifyCryptoData.HandlerFunctions
  def prepare_purge() do
    Task.Supervisor.start_child(Purgehelper.Supervisor, &make_purge/0)
  end

  defp make_purge() do
    HandlerFunctions.get_deletables_transaction()
    |> Stream.chunk_every(10)
    |> Stream.map(&handle_purge/1)
    |> Stream.each(fn task -> Task.await(task) end)
    |> Stream.run()
  end

  defp handle_purge(terms) do
    Task.async(fn ->
      Poolex.run(
        :purge,
        fn pid ->
          GenServer.call(pid, {:purge, terms})
        end,
        checkout_timeout: 60000
      )
    end)
  end
end
