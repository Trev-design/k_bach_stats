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
          resolve(useQuery(GET_ACCOUNT))
        }, 200)
      })
    }
    
    const {result, loading, error} = getAccountDetails()

    console.log(result)
    console.log(loading)
    console.log(error)

    return {
      data: result,
      loading,
      error
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