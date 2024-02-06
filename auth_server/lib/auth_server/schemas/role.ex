defmodule AuthServer.Schemas.Role do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "roles" do
    field :verificated, :boolean
    belongs_to :user, AuthServer.Schemas.User
    timestamps()
  end

  def changeset(role, attrs) do
    role
    |> cast(attrs, [:user_id, :verificated])
    |> validate_required([:user_id, :verificated])
  end
end
