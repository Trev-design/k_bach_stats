defmodule AuthServer.RolesFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `AuthServer.Roles` context.
  """

  @doc """
  Generate a role.
  """
  def role_fixture(attrs \\ %{}) do
    {:ok, role} =
      attrs
      |> Enum.into(%{
        trait: "some trait",
        verified: true
      })
      |> AuthServer.Roles.create_role()

    role
  end
end
