<template>
  <section class="profile-page-container">
    <section class="profile-details-container">
      <div class="profile-image-container">
        <img :src="imageFilePath" alt="popel" class="profile-image">
        <h3 class="name"> {{ name }} </h3>
      </div>
      <div class="profile-details">
        <p class="detail-field"> {{ name }} </p>
        <p class="detail-field"> {{ email }} </p>
      </div>
    </section>
    <div class="profile-edit">
      <router-link class="edit-router-link">settings</router-link>
      <router-link class="edit-router-link">edit profile</router-link>
    </div>
    <div class="bio-container">
      <p class="bio"></p>
    </div>
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

  async created() {
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
  width: 85vw;
}

.profile-details-container {
  display: flex;
  flex-direction: row;
  justify-content: space-arround;
  height: 200px;
}

.profile-image-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: left;
  width: 40%;
  padding-left: 200px;
}

.profile-image {
  width: 80px;
  height: 80px;
}

.profile-details {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: left;
}

.name {
  font-size: 1.1rem;
  color: rgb(170, 225, 205);
}

.profile-edit {
  display: flex;
  justify-content: right;
  align-items: center;
}

.edit-router-link {
  text-decoration: none;
  padding: 0 0.75rem;
  font-size: 1.1rem;
  color: rgb(100, 175, 125);
  transition: all 0.3s ease;
}

.edit-router-link:hover {
  color: rgb(215, 125, 225);
}

.detail-field {
  font-size: 1.1rem;
  color: rgb(170, 225, 205);
  padding: 0.6rem 0;
}

</style>