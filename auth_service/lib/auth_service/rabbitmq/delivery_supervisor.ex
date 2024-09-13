defmodule AuthService.Rabbitmq.DeliverySupervisor do
  use Supervisor

  def start_link([]), do: Supervisor.start_link(__MODULE__, [], name: __MODULE__)

  def init([]) do
    children = [
      {Poolex, [
        pool_id: :rabbit_delivery,
        worker_module: AuthService.Rabbitmq.Publisher,
        workers_count: 5,
        max_overflow: 1,
      ]}
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
