defmodule EmailTest do
  use ExUnit.Case

  test "sent emails" do
    email_delivery = AuthServer.JobHandler.EmailJob.deliver("ass", "ass@ass.ass")

    assert email_delivery == :ok
  end
end
