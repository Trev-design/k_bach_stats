defmodule AuthServer.SessionHandler do
  require Ecto.Query
  alias AuthServer.{Repo, Schemas.Account, Schemas.User, Schemas.Role}

  def register(name, email, password) do
    with {:ok, password_hash}        <- try_make_password_hash(password),
         {:ok, %Account{} = account} <- create_account(%{email: email, password_hash: password_hash}),
         {:ok, %User{} = user}       <- create_user(account, %{name: name})
    do
      {200, %{name: user.name, id: user.id}}
    else
      {:error, %Ecto.Changeset{} = changeset} -> {403, %{message: errors(changeset.errors)}}
      {:error, reason}                        -> {403, %{message: reason}}
    end
  end

  def signin(email, password) do
    with %Account{} = account <- Repo.get_by(Account, email: email),
         true                 <- Argon2.verify_pass(password, account.password_hash)
    do
      {:ok, Repo.preload(account, :user)}
    else
      nil               -> {:error, "could not find an account wit this email"}
      false             -> {:error, "password not correct"}
    end
  end

  def check_verification(verification, cookie) do
    with {:ok, payload}         <- Jason.decode(cookie),
         {:ok, user}            <- get_user(payload["id"]),
         true                   <- payload["verification"] == verification,
         {:ok, %Role{} = _role} <- create_role(user, %{verificated: true})
    do
      {:ok, user}
    else
      true                                    -> {:error, "false verification code"}
      {:error, %Ecto.Changeset{} = changeset} -> {:error, errors(changeset.errors)}
      error                                   -> error
    end
  end

  def get_user(id) do
    case User |> Repo.get(id) do
      %User{} = user -> {:ok, user}
      nil            -> {:error, "no user with this id"}
    end
  end

  defp try_make_password_hash(password) do
    case Regex.match?(~r/^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!?&%$§#@€]).{10,}$/, password) do
      true  -> {:ok, Argon2.hash_pwd_salt(password)}
      false -> {:error, "invalid password"}
    end
  end

  defp create_account(attrs) do
    %Account{}
    |> Account.changeset(attrs)
    |> Repo.insert()
  end

  defp create_user(%Account{} = account, attrs) do
    account
    |> Ecto.build_assoc(:user)
    |> User.changeset(attrs)
    |> Repo.insert()
  end

  defp errors(error_list) do
    error_list
    |> Keyword.values()
    |> Enum.map(fn error -> elem(error, 0) end)
  end

  defp create_role(%User{} = user, attrs) do
    user
    |> Ecto.build_assoc(:role)
    |> Role.changeset(attrs)
    |> Repo.insert()
  end
end
