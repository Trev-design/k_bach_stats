<template>
  <div class="navbar-container">
    <ul class="nav-items">
      <li><router-link class="nav-link">{{ localStorage.getItem('guest') }}</router-link></li>
      <li><router-link to="/signin" class="nav-link" @click="handleLogout()">Logout</router-link></li>
    </ul>
  </div>
</template>


<script>
export default {
  name: 'HomeNavbar',
  props: {
    handleLogout: Function
  } 
}
</script>


<style scoped>
.navbar-container {
  position: fixed;
  top: 0;
  z-index: 2;
  width: 100%;
  display: flex;
  justify-content: end;
  align-items: center;
}

.nav-items {
  display: flex;
  flex-direction: row;

  li {
    list-style: none;
  }
}

.nav-link {
  margin: 0 1.2rem 0 0;
  text-decoration: none;
}
</style>