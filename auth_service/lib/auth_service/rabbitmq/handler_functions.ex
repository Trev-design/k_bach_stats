defmodule AuthService.Rabbitmq.HandlerFunctions do
  use AMQP

  @spec setup_connections(dsn :: String.t()) :: AMQP.Channel.t()
  def setup_connections(dsn) do
    {:ok, conn} = Connection.open(dsn)
    {:ok, chan} = Channel.open(conn)
    chan
  end

  @spec declare_exchange(channel :: AMQP.Channel.t(), exchange :: String.t(), durable :: boolean()) :: AMQP.Channel.t()
  def declare_exchange(channel, exchange, durable) do
    :ok = Exchange.declare(
      channel,
      exchange,
      :direct,
      durable: durable,
      auto_delete: false,
      internal: false,
      no_wait: false
    )

    channel
  end

  @spec declare_queue(channel :: AMQP.Channel.t(), queue :: String.t(), durable :: boolean(), auto_delete :: boolean()) :: AMQP.Channel.t()
  def declare_queue(channel, queue, durable, auto_delete) do
    case Queue.declare(
      channel,
      queue,
      durable: durable,
      auto_delete: auto_delete,
      exclusive: false,
      nowait: false)
    do
      :ok      -> channel
      {:ok, _} -> channel
      _invalid -> raise "invalid queue declare"
    end
  end

  @spec bind_queue(channel :: AMQP.Channel.t(), exchange :: String.t(), routing_key :: String.t(), queue :: String.t()) :: AMQP.Channel.t()
  def bind_queue(channel, exchange, routing_key, queue) do
    :ok = Queue.bind(
      channel,
      queue,
      exchange,
      routing_key: routing_key,
      nowait: false
    )

    channel
  end

  @spec publish(channel :: AMQP.Channel.t(), exchange :: String.t(), routing_key :: String.t(), payload :: String.t()) :: :ok | {:error, reason :: any()}
  def publish(channel, exchange, routing_key, payload) do
    Basic.publish(
      channel,
      exchange,
      routing_key,
      payload,
      persistent: true,
      mandatory: true,
      correlation_id: Uniq.UUID.uuid4(),
      content_type: "application/json"
    )
  end
end
