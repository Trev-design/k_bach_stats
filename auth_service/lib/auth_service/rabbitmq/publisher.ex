defmodule AuthService.Rabbitmq.Publisher do
  alias AuthService.Rabbitmq.ConnectionHandler
  alias AuthService.Rabbitmq.HandlerFunctions
  alias AuthService.Accounts.Account

  use GenServer

  require Logger

  def start_link(), do: GenServer.start_link(__MODULE__, [])

  def init([]), do: {:ok, nil, {:continue, :declare_channel}}

  def handle_info({:set_to_nil, pid}, _state), do: {:noreply, nil, {:continue, {:terminated, pid}}}

  def handle_continue(:declare_channel, _state) do
    {:noreply, ConnectionHandler.declare_channel(self())}
  end

  def handle_continue({:terminated, pid}, state) do
    send(pid, {:ok, :success})
    {:noreply, state}
  end

  def handle_call({:verify_email, username, id, email, validation}, _from, channel) do
    result =
      HandlerFunctions.publish(
        channel,
        "verify",
        "send_verify_email",
        Jason.encode!(%{
          kind: "verify",
          email: email,
          verify: validation,
          name: username,
          user_id: id
        }))

    case result do
      :ok -> {:reply, :published, channel}
      invalid ->
        # TODO: implement background error handler
        {:reply, invalid, channel}
    end
  end

  def handle_call({:forgotten_password, username, id, email, validation}, _from, channel) do
    result =
        HandlerFunctions.publish(
        channel,
        "verify",
        "send_forgotten_password_email",
        Jason.encode!(%{
          kind: "forgotten_password",
          email: email,
          verify: validation,
          name: username,
          user_id: id}))

    case result do
      :ok -> {:reply, :published, channel}
      invalid ->
        # TODO: implement background error handler
        {:reply, invalid, channel}
    end
  end

  def handle_call({:session, username, account, id, abo}, _from, channel) do
    result =
      HandlerFunctions.publish(
        channel,
        "session",
        "send_session_credentials",
        Jason.encode!(%{
          name: username,
          account: account,
          id: id,
          abo_type: abo}))

    case result do
      :ok -> {:reply, :published, channel}
      invalid ->
        # TODO: implement background error handler
        {:reply, invalid, channel}
    end
  end

  def handle_call({:enroll_account, %Account{} = account, id, abo}, _from, channel) do
    user_account_payload = Jason.encode!(%{entity: account.id, username: account.user.id, email: account.email, abo_type: abo})
    session_payload = Jason.encode!(%{name: account.user.name, account: account.id, id: id})

    account_result =
      HandlerFunctions.publish(
        channel,
        "account",
        "add_account_request",
        user_account_payload
      )

    Logger.info("set account successful")

    session_result =
      HandlerFunctions.publish(
        channel,
        "session",
        "send_session_credentials",
        session_payload
      )

    Logger.info("set session successful")

    with :ok <- account_result,
         :ok <- session_result
    do
      {:reply, :enrolled_user, channel}
    else
      invalid ->
        # TODO: implement background error handler
        {:reply, invalid, channel}
    end
  end

  def terminate(reason, channel) do
    Logger.info("i am terminating here")
    Logger.info(reason)
    case channel do
      nil   -> {:stop, reason}
      _chan ->
        ConnectionHandler.stop_channel(channel)
        {:stop, reason}
    end
  end
end
