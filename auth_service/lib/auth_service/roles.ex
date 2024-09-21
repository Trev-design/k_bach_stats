defmodule AuthService.Roles do
  @moduledoc """
  The Roles context.
  """

  import Ecto.Query, warn: false
  alias AuthService.Repo

  alias AuthService.Roles.Role


  def list_roles do
    Repo.all(Role)
  end

  def get_role!(id), do: Repo.get!(Role, id)

  def create_role(account, attrs \\ %{}) do
    account
    |> Ecto.build_assoc(:role)
    |> Role.changeset(attrs)
    |> Repo.insert()
  end

  def update_role(%Role{} = role, attrs) do
    role
    |> Ecto.Changeset.change(attrs)
    |> Repo.update()
  end

  def delete_role(%Role{} = role) do
    Repo.delete(role)
  end

  def change_role(%Role{} = role, attrs \\ %{}) do
    Role.changeset(role, attrs)
  end
end
