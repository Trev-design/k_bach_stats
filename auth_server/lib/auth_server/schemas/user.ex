defmodule AuthServer.Schemas.User do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "users" do
    field :name, :string
    has_one :role, AuthServer.Schemas.Role
    belongs_to :account, AuthServer.Schemas.Account

    timestamps()
  end

  def changeset(user, attrs) do
    user
    |> cast(attrs, [:account_id, :name])
    |> validate_required([:account_id])
  end
end
