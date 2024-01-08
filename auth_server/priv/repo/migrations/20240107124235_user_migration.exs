defmodule AuthServer.Repo.Migrations.UserMigration do
  use Ecto.Migration

  def change do
    create table(:users, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :name, :string
      add :account_id, references(:accounts, on_delete: :delete_all, type: :binary_id)

      timestamps()
    end
  end
end
