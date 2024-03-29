defmodule AuthServer.SessionHandler do
  require Ecto.Query
  require Logger

  alias AuthServer.{
    Repo,
    Schemas.Account,
    Schemas.User,
    Schemas.Role
  }

  def register(name, email, password) do
    with {:ok, password_hash}        <- try_make_password_hash(password),
         {:ok, "email valid"}        <- check_email(email),
         {:ok, %Account{} = account} <- create_account(%{email: email, password_hash: password_hash}),
         {:ok, %User{} = user}       <- create_user(account, %{name: name}),
         {:ok, %Role{}}              <- create_role(user, %{verified: false})
    do
      {200, %{name: user.name, id: user.id}}
    else
      {:error, %Ecto.Changeset{} = changeset} -> {403, %{message: errors(changeset.errors)}}
      {:error, reason}                        -> {403, %{message: reason}}
    end
  end

  def signin(email, password) do
    with %Account{} = account <- get_by_email(email),
         true                 <- Argon2.verify_pass(password, account.password_hash)
    do
      if account.user.role.verified, do: {:ok, account}, else: {:error, "user is not verified"}
    else
      nil   -> {:error, "could not find an account wit this email"}
      false -> {:error, "password not correct"}
    end
  end

  def check_verification(verification, cookie) do
    Logger.info(inspect(cookie))
    with {:ok, payload} <- Jason.decode(cookie),
         {:ok, user}    <- get_user(payload["id"]),
         {:ok, integer} <- verification_integer(verification),
         {:ok, %Role{}} <- update_role(user.role, %{verified: true})
    do
      if integer == payload["verify"], do: {:ok, user}, else: {:error, "verification does not match"}
    else
      {:error, %Ecto.Changeset{} = changeset} -> {:error, errors(changeset.errors)}
      error                                   -> error
    end
  end

  def check_email_verification(verification, cookie) do
    with {:ok, payload} <- Jason.decode(cookie),
         {:ok, account} <- get_account(payload["id"]),
         {:ok, _integer} <- verification_integer(verification)
    do
      {:ok, account}
    else
      error -> error
    end
  end

  def get_account(id) do
    case Repo.get(Account, id) do
      %Account{} = account -> {:ok, account}
      nil                  -> {:error, "no account with this id"}
    end
  end

  def get_user(id) do
    case User |> Repo.get(id) |> Repo.preload(:role) do
      %User{} = user -> {:ok, user}
      nil            -> {:error, "no user with this id"}
    end
  end

  def get_by_email(email), do: Repo.get_by(Account, email: email) |> Repo.preload(user: :role)

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

  defp update_role(old_role, new_attrs) do
    old_role
    |> Role.changeset(new_attrs)
    |> Repo.update()
  end

  def update_password(old_account, new_password) do
    with {:ok, password}             <- try_make_password_hash(new_password),
         {:ok, %Account{} = account} <- update_password_in_account(old_account, password)
    do
      {:ok, account}
    else
      {:error, %Ecto.Changeset{} = changeset} -> {:error, errors(changeset.errors)}
      error                                   -> error
    end
  end

  defp check_email(email) do
    case Regex.match?(~r/^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[a-z]{2,4}$/, email) do
      true  -> {:ok, "email valid"}
      false -> {:error, "invalid email"}
    end
  end

  defp verification_integer(verification) do
    case Integer.parse(verification) do
      {integer, ""} -> {:ok, integer}
      _invalid      -> {:error, "verification not just an integer"}
    end
  end

  defp update_password_in_account(old_account, new_password) do
    old_account
    |> Account.changeset(%{password: new_password})
    |> Repo.update()
  end
end
