defmodule AuthServer.Accounts do
  @moduledoc """
  The Accounts context.
  """

  import Ecto.Query, warn: false
  alias AuthServer.Repo

  alias AuthServer.Accounts.Account


  def list_accounts do
    Repo.all(Account)
  end

  def get_account(id), do: Repo.get(Account, id)

  def get_full_account(id), do: Account |> Repo.get(id) |> Repo.preload(user: :role)
  def get_account_by_email(email), do: Repo.get_by(Account, email: email) |> Repo.preload(user: :role)

  def create_account(attrs \\ %{}) do
    %Account{}
    |> Account.changeset(attrs)
    |> Repo.insert()
  end

  def change_password(%Account{} = account, new_password_hash) do
    account
    |> Ecto.Changeset.change(%{password_hash: new_password_hash})
    |> Repo.update()
  end

  def update_account(%Account{} = account, attrs) do
    account
    |> Account.changeset(attrs)
    |> Repo.update()
  end

  def delete_account(%Account{} = account) do
    Repo.delete(account)
  end

  def change_account(%Account{} = account, attrs \\ %{}) do
    Account.changeset(account, attrs)
  end
end
