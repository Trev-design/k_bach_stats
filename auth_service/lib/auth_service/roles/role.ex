defmodule AuthService.Roles.Role do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "roles" do
    field :verified, :boolean, default: false
    field :abo_type, :string, default: "NOT_VERIFIED"

    belongs_to :account, AuthService.Accounts.Account

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(role, attrs) do
    role
    |> cast(attrs, [:account_id])
    |> validate_required([:account_id])
  end
end
