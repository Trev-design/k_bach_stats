defmodule AuthServer.SessionHandler do

  def verify_code() do
    for _x <- 1..7 do
      :rand.uniform(9) + 48
    end
    |> List.to_integer()
  end

  def send_verify_email(name, email, verify_code, route) do
    case Finch.build(
          :post, "http://127.0.0.1:8080/#{route}",
          [{"Content-Type", "application/json"}],
          Jason.encode!(
            %{
              to: email,
              subject: name,
              data: "Hello #{name} your verify code is #{verify_code}"
              }))
         |> Finch.request(EmailRequest)
    do
      {:ok, %Finch.Response{}} -> :ok
      {:error, _exception} -> :error
    end
  end

  def errors(error_list) do
    error_list
    |> Keyword.values()
    |> Enum.map(fn error_tuple -> elem(error_tuple, 0) end)
  end

  def check_password(password, password_hash) do
    case Argon2.verify_pass(password, password_hash) do
      true  -> :ok
      false -> {:error, "invalid password"}
    end
  end
end
