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
import { GET_ACCOUNT } from '../queries'

export default {
  name: 'Profile',

  data() {
    return {
      name: '',
      email: '',
      imageFilePath: '',
      bio: '',
      workspaces: []
    }
  },

  async mounted() {
    const id = localStorage.getItem('initialUser')
      const { data } = await this.$apollo.query({
        query: GET_ACCOUNT,
        variables: {userID: id}
      })

      this.name = data.getUser.profile.contact.name
      this.email = data.getUser.profile.contact.email
      this.imageFilePath = data.getUser.profile.contact.imageFilePath
      this.bio = data.getUser.profile.bio
      
      if (data.getUser.workspaces.length > 0) {
        this.workspaces.push(...data.getUser.workspaces)
      }
  }
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