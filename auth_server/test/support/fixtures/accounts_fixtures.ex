defmodule AuthServer.AccountsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `AuthServer.Accounts` context.
  """

  @doc """
  Generate a account.
  """
  def account_fixture(attrs \\ %{}) do
    {:ok, account} =
      attrs
      |> Enum.into(%{
        email: "some email",
        password_hash: "some password_hash"
      })
      |> AuthServer.Accounts.create_account()

    account
  end
end
