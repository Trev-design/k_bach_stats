defmodule AuthService.Rabbitmq.RabbitInstance do
  use Supervisor

  def start_link(props), do: Supervisor.start_link(__MODULE__, props, name: __MODULE__)

  def init(props) do
    children = [
      {AuthService.Rabbitmq.ConnectionHandler, props},
      {AuthService.Rabbitmq.ChannelHandler, []},
      {AuthService.Rabbitmq.DeliverySupervisor, []}
    ]

    Supervisor.init(children, strategy: :rest_for_one)
  end
end
