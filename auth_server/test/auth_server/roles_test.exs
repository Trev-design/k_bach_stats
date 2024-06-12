defmodule AuthServer.RolesTest do
  use AuthServer.DataCase

  alias AuthServer.Roles

  describe "roles" do
    alias AuthServer.Roles.Role

    import AuthServer.RolesFixtures

    @invalid_attrs %{role: nil, trait: nil}

    test "list_roles/0 returns all roles" do
      role = role_fixture()
      assert Roles.list_roles() == [role]
    end

    test "get_role!/1 returns the role with given id" do
      role = role_fixture()
      assert Roles.get_role!(role.id) == role
    end

    test "create_role/1 with valid data creates a role" do
      valid_attrs = %{role: "some role", trait: "some trait"}

      assert {:ok, %Role{} = role} = Roles.create_role(valid_attrs)
      assert role.role == "some role"
      assert role.trait == "some trait"
    end

    test "create_role/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Roles.create_role(@invalid_attrs)
    end

    test "update_role/2 with valid data updates the role" do
      role = role_fixture()
      update_attrs = %{role: "some updated role", trait: "some updated trait"}

      assert {:ok, %Role{} = role} = Roles.update_role(role, update_attrs)
      assert role.role == "some updated role"
      assert role.trait == "some updated trait"
    end

    test "update_role/2 with invalid data returns error changeset" do
      role = role_fixture()
      assert {:error, %Ecto.Changeset{}} = Roles.update_role(role, @invalid_attrs)
      assert role == Roles.get_role!(role.id)
    end

    test "delete_role/1 deletes the role" do
      role = role_fixture()
      assert {:ok, %Role{}} = Roles.delete_role(role)
      assert_raise Ecto.NoResultsError, fn -> Roles.get_role!(role.id) end
    end

    test "change_role/1 returns a role changeset" do
      role = role_fixture()
      assert %Ecto.Changeset{} = Roles.change_role(role)
    end
  end

  describe "roles" do
    alias AuthServer.Roles.Role

    import AuthServer.RolesFixtures

    @invalid_attrs %{verified: nil, trait: nil}

    test "list_roles/0 returns all roles" do
      role = role_fixture()
      assert Roles.list_roles() == [role]
    end

    test "get_role!/1 returns the role with given id" do
      role = role_fixture()
      assert Roles.get_role!(role.id) == role
    end

    test "create_role/1 with valid data creates a role" do
      valid_attrs = %{verified: true, trait: "some trait"}

      assert {:ok, %Role{} = role} = Roles.create_role(valid_attrs)
      assert role.verified == true
      assert role.trait == "some trait"
    end

    test "create_role/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Roles.create_role(@invalid_attrs)
    end

    test "update_role/2 with valid data updates the role" do
      role = role_fixture()
      update_attrs = %{verified: false, trait: "some updated trait"}

      assert {:ok, %Role{} = role} = Roles.update_role(role, update_attrs)
      assert role.verified == false
      assert role.trait == "some updated trait"
    end

    test "update_role/2 with invalid data returns error changeset" do
      role = role_fixture()
      assert {:error, %Ecto.Changeset{}} = Roles.update_role(role, @invalid_attrs)
      assert role == Roles.get_role!(role.id)
    end

    test "delete_role/1 deletes the role" do
      role = role_fixture()
      assert {:ok, %Role{}} = Roles.delete_role(role)
      assert_raise Ecto.NoResultsError, fn -> Roles.get_role!(role.id) end
    end

    test "change_role/1 returns a role changeset" do
      role = role_fixture()
      assert %Ecto.Changeset{} = Roles.change_role(role)
    end
  end
end
