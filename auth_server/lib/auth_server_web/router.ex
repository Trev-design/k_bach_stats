defmodule AuthServerWeb.Router do
  use AuthServerWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
    plug :fetch_session
  end

  pipeline :verify do
    plug AuthServerWeb.VerifyPlug
  end

  pipeline :auth do
    plug AuthServerWeb.AuthPlug
  end

  scope "/account", AuthServerWeb do
    pipe_through :api
    post "/signup", AccountController, :signup
    post "/signin", AccountController, :signin
    post "/new_verify", ChangeController, :new_verify
    post "/forgotten_password", ChangeController, :forgotten_password
  end

  scope "/account/verify", AuthServerWeb do
    pipe_through [:api, :verify]
    post "/user", VerifyController, :verify
    post "/new_password", VerifyController, :change_password
  end

  scope "/account/session", AuthServerWeb do
    pipe_through [:api, :auth]
    get "/refresh", SessionController, :refresh_token
    get "/signout", SessionController, :logout
  end
end
