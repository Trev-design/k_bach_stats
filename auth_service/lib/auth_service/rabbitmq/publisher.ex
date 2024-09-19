defmodule AuthService.Rabbitmq.Publisher do
  alias AuthService.Accounts.Account
  alias AuthService.Rabbitmq.ChannelHandler
  alias AuthService.Rabbitmq.HandlerFunctions
  use GenServer

  require Logger

  def start_link(), do: GenServer.start_link(__MODULE__, [])

  def init([]), do: {:ok, nil, {:continue, :declare_channel}}

  def handle_continue(:declare_channel, state) do
    send(:rmq_connection, {:declare_channel, self()})
    {:noreply, state}
  end

  def handle_info({:get_connection, connection}, _state) do
    {:ok, channel} = HandlerFunctions.setup_channel(connection)
    Logger.info("setup channel")
    ChannelHandler.set_channel(channel)
    {:noreply, channel}
  end

  def handle_call({:enroll_account, %Account{} = account, id}, _from, channel) do
    user_account_payload = Jason.encode!(%{entity: account.id, username: account.user.id, email: account.email})
    session_payload = Jason.encode!(%{name: account.user.name, account: account.id, id: id})

    account_result =
      HandlerFunctions.publish(
        channel,
        "account",
        "add_account_request",
        user_account_payload
      )

    session_result =
      HandlerFunctions.publish(
        channel,
        "session",
        "send_session_credentials",
        session_payload
      )

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

  def handle_call({:session, username, account, id}, _from, channel) do
    result =
      HandlerFunctions.publish(
        channel,
        "session",
        "send_session_credentials",
        Jason.encode!(%{
          name: username,
          account: account,
          id: id}))

    case result do
      :ok -> {:reply, :published, channel}
      invalid ->
        # TODO: implement background error handler
        {:reply, invalid, channel}
    end
  end

  def terminate(reason, channel) do
    send(:rmq_connection, {:stop_channel, channel})
    {:stop, reason}
  end
end
