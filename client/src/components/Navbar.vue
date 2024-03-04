<template>
  <div class="navbar-container">
    <ul class="nav-items">
      <li><router-link class="navlink" to="/register">Register</router-link></li>
      <li><router-link class="navlink" to="/signin">Signin</router-link></li>
    </ul>
  </div>
</template>


<script>
export default {
  name: 'Navbar',
  data: () => ({
    navbarData: [{text: "signup"}]
  }) 
}
</script>


<style>
.navbar-container{
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

.navlink {
  margin: 0 1.2rem 0 0;
  text-decoration: none;
}
</style>