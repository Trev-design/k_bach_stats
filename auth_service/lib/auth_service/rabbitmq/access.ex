defmodule AuthService.Rabbitmq.Access do

  def publish_verify_message(validation, username, email, account) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:verify_email, username, account, email, validation})
    end))
  end

  def publish_forgotten_password_message(validation, username, email, account) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:forgotten_password, username, account, email, validation})
    end))
  end

  def publish_session_message(username, account, session) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:session, username, account, session})
    end))
  end

  defp handle_publish(fun) do
    Task.async(fn ->
      Poolex.run(
        :rabbit_delivery,
        fun,
        checkout_timeout: 10000
      )
    end)
  end
end
