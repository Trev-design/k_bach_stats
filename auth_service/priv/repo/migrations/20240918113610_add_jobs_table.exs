defmodule AuthService.Repo.Migrations.AddJobsTable do
  use Ecto.Migration

  def up, do: Oban.Migrations.up()
  def down, do: Oban.Migration.down(version: 1)
end
