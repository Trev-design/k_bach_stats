defmodule AuthServiceWeb.Router do
  use AuthServiceWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
    plug :fetch_session
  end

  pipeline :verify do
    plug AuthServiceWeb.VerifyPlug
  end

  pipeline :auth do
    plug AuthServiceWeb.AuthPlug
  end

  scope "/account", AuthServiceWeb do
    pipe_through :api

    post "/signup", AccountController, :signup
    post "/signin", AccountController, :signin
    post "/new_verify", AccountController, :request_new_verify
  end

  scope "/verify", AuthServiceWeb do
    pipe_through [:api, :verify]

    post "/account", VerifyController, :verify
  end

  scope "/forgotten_password", AuthServiceWeb do
    pipe_through [:api, :verify]

    post "/", VerifyController, :forgotten_password
  end

  scope "/session", AuthServiceWeb do
    pipe_through [:api, :auth]

    get "/refresh", SessionController, :refresh_session
    get "/signout", SessionController, :signout
  end
end
