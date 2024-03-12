<template>
  <div class="navbar-container">
    <h1 class="brand">KBach</h1>
    <a href="#about" class="navlink">About</a>
    <a href="#features" class="navlink">Features</a>
    <a href="#testimonials" class="navlink">Testimonials</a>
    <a href="#pricing" class="navlink">Pricing</a>
    <ul class="nav-items">
      <li><router-link class="navlink" to="/register">Register</router-link></li>
      <li><router-link class="navlink" to="/signin">Signin</router-link></li>
    </ul>
  </div>
</template>


<script>
export default {
  name: 'LandingNavbar'
}
</script>


<style scoped>
.navbar-container{
  position: fixed;
  top: 0;
  z-index: 2;
  width: 100%;
  display: flex;
  justify-content: end;
  align-items: center;
  background-color: rgb(3, 6, 18);
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
  font-size: 1.1rem;
  color: rgb(150, 165, 225);
  position: relative;
  &::before {
    content: '';
    width: 100%;
    bottom: -4px;
    height: 2px;
    background-color: rgb(110, 170, 250);
    position: absolute;
    transform: scale(0, 1);
    transition: all 0.3s ease;
  }

  &:hover {&::before {transform: scale(1, 1);}}
}

.brand {
  background: rgb(17,158,186);
  background: linear-gradient(90deg, rgba(17,158,186,1) 16%, rgba(43,43,222,1) 78%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-right: auto;
  margin-left: 1rem;
}
</style>