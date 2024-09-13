defmodule AuthService.Rabbitmq.HandlerFunctions do
  use AMQP

  def setup_connections(dsn), do: Connection.open(dsn)

  def setup_channel(connection), do: Channel.open(connection)

  # @spec declare_exchange(channel :: AMQP.Channel.t(), exchange :: String.t(), durable :: boolean()) :: AMQP.Channel.t()
  # def declare_exchange(channel, exchange, durable) do
  #   :ok = Exchange.declare(
  #     channel,
  #     exchange,
  #     :direct,
  #     durable: durable,
  #     auto_delete: false,
  #     internal: false,
  #     no_wait: false
  #   )

  #   channel
  # end

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

  def close_channel(channel), do: Channel.close(channel)
  def close_connection(connection), do: Connection.close(connection)
end
