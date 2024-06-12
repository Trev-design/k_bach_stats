defmodule AuthServer.Roles.Role do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "roles" do
    field :verified, :boolean, default: false
    field :trait, :string
    belongs_to :user, AuthServer.Users.User

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(role, attrs) do
    role
    |> cast(attrs, [:verified, :trait, :user_id])
    |> validate_required([:verified, :trait, :user_id])
  end
end
