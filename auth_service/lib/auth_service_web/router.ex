defmodule AuthServiceWeb.Router do
  use AuthServiceWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
    plug :fetch_session
  end

  pipeline :verify do
    plug AuthServiceWeb.VerifyPlug
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
end
