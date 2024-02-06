defmodule AuthServer.Routers.RouterHelpers do
  alias AuthServer.Jwt

  def add_claims_and_generate(user_id, user_name) do
    Jwt.generate_token_pair(
      %{"id" => user_id,
        "name" => user_name,
        "exp" => Joken.current_time() + (15 * 60)},
      %{"id" => user_id,
        "exp" => Joken.current_time()+ (24 * 60 * 60)
      }
    )
  end
end
