<script>
export default {
  name: 'App',
  data: () => ({
    firstName: 'Mehners',
    lastName: 'Moos'
  }),
  computed: {
    fullName() {
      return this.firstName + " " + this.lastName
    }
  }
}
</script>


<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>


<style>

</style>
