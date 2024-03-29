defmodule AuthServer.Schemas.Account do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "accounts" do
    field :email, :string
    field :password_hash, :string
    has_one :user, AuthServer.Schemas.User

    timestamps()
  end

  def changeset(account, attrs) do
    account
    |> cast(attrs, [:email, :password_hash])
    |> validate_required([:email, :password_hash])
    |> unique_constraint(:email)
  end
end
