defmodule AuthServiceWeb.Router do
  use AuthServiceWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
    plug :fetch_session
  end

  scope "/account", AuthServiceWeb do
    pipe_through :api

    post "/signup", AccountController, :signup
  end
end
