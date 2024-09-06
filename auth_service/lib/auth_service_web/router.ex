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
  end

  scope "/verify", AuthServiceWeb do
    pipe_through [:api, :verify]

    post "/account", VerifyController, :verify
  end
end
