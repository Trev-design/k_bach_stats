<template>
  <section class="profile-page-container">
    <section class="profile-details-container">
      <div class="profile-image-container">
        <img :src="imageFilePath" alt="popel" class="profile-image">
        <h3 class="name"> {{ name }} </h3>
      </div>
      <div class="profile-details">
        <p class="name-field"> {{ name }} </p>
        <p class="email-f"> {{ email }} </p>
      </div>
    </section>
    <div class="profile-edit">
      <router-link class="edit-router-link">settings</router-link>
      <router-link class="edit-router-link">edit profile</router-link>
    </div>
    <section class="workspace">

    </section>
  </section>
</template>

<script>
import { useQuery } from '@vue/apollo-composable';
import { GET_ACCOUNT } from '../queries'
export default {
  name: 'Profile',

  data() {
    return {
      name: '',
      email: '',
      imageFilePath: '',
      bio: '',
      workspaces: [],
    }
  },

  methods: {
    fetchData() {
      const id = localStorage.getItem('initialUser')
      const {result, loading, error} = useQuery(GET_ACCOUNT, {userID: id})
      console.log(loading)
      console.log(result.value)
      console.log(error)

      this.name = result.value.getUser.profile.contact.name
      this.email = result.value.getUser.profile.contact.email
      this.imageFilePath = result.value.getUser.profile.contact.imageFilePath
      this.bio = result.value.getUser.profile.bio

      console.log(result.value.getUser.profile.contact.name)
      console.log(result.value.getUser.profile.contact.email)

      if (result.value.getUser.workspaces.length > 0) {
        this.workspaces.push(...result.value.getUser.workspaces)
      }
      console.log(`feched data for ${result.value.getUser.entity}`)
    }
  },

  mounted() {
    this.fetchData()
  },
}
</script>

<style scoped>
.profile-page-container {
  display: flex;
  flex-direction: column;
  width: 100vw;
}

.profile-details-container {
  display: flex;
  flex-direction: row;
  width: 100vw;
  height: 200px;
  justify-content: space-around
}

.profile-image-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>