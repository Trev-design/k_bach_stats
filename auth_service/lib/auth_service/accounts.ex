defmodule AuthService.Accounts do
  @moduledoc """
  The Accounts context.
  """

  import Ecto.Query, warn: false
  alias AuthService.Repo

  alias AuthService.Accounts.Account
  def list_accounts do
    Repo.all(Account)
  end

  def get_full_account_by_email(email), do: Account |> Repo.get_by(email: email) |> Repo.preload([:user, :role])
  def get_full_account(id), do: Account |> Repo.get(id) |> Repo.preload([:user, :role])

  def get_account(id), do: Repo.get(Account, id)
  def get_account_by_email(email), do: Account |> Repo.get_by(email: email)

  def create_account(attrs \\ %{}) do
    %Account{}
    |> Account.changeset(attrs)
    |> Repo.insert()
  end

  def update_account(%Account{} = account, attrs) do
    account
    |> Account.new_password_changeset(attrs)
    |> Repo.update()
  end

  def delete_account(%Account{} = account) do
    Repo.delete(account)
  end

  def change_account(%Account{} = account, attrs \\ %{}) do
    Account.changeset(account, attrs)
  end
end
