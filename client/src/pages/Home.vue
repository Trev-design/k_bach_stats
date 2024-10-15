<template>
  <section class="home-container">
    <section class="home-account-content-container">

    </section>
  </section>
</template>

<script>
import { useQuery } from '@vue/apollo-composable';
import { GET_ACCOUNT } from '../queries';


export default {
  name: 'Home',

  setup() {
    const getAccountDetails = () => {
      return new Promise((resolve) => {
        setTimeout(() => {
          const {result, loading, error} = useQuery(GET_ACCOUNT, {})
          console.log(result)
          console.log(loading)
          console.log(error)
          resolve(result)
        }, 2000)
      })
    }

    getAccountDetails()

    return {
      data: null
    }
  }
}
</script>

<style scoped>
.home-container {
  width: 100%;
  height: 100%;
}

.home-account-content-container {
  position: relative;
  left: 2rem;
  width: 85vw;
  height: 100vh;
  margin: 7rem;
  border: 1px solid white;
}
</style>