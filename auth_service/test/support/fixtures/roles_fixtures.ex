defmodule AuthService.RolesFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `AuthService.Roles` context.
  """

  @doc """
  Generate a role.
  """
  def role_fixture(attrs \\ %{}) do
    {:ok, role} =
      attrs
      |> Enum.into(%{
        abo_type: "some abo_type",
        verified: true
      })
      |> AuthService.Roles.create_role()

    role
  end
end
