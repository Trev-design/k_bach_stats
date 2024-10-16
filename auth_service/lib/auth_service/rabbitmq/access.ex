defmodule AuthService.Rabbitmq.Access do
  alias AuthService.Accounts.Account
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

  def publish_session_message(username, account, session, abo) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:session, username, account, session, abo})
    end))
  end

  def publish_enroll_user(%Account{} = account, session, role) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:enroll_account, account, session, role})
    end))
  end

  def publish_remove_session(%Account{} = account, session) do
    Task.await(handle_publish(fn pid ->
      GenServer.call(pid, {:remove_session, account, session})
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
