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

    <section class="person-details-container">
      <div class="expierence-container">
        <div 
          v-for="experience in experiences"
          :key="experience.id" 
          class="experience"
        >
          <p class="experience-description"> {{ experience.experience }} </p>
          <Icon 
            v-for="star in 5"
            :key="star" 
            :class="`experience-star-rating-container ${star <= experience.rating ? 'filled' : ''}`" 
            icon="material-symbols:star-outline-rounded"
          />
        </div>

        <div class="add-experience-button-container">
          <button 
            class="add-experience-button"
            @click="handleAddExperience()"
          >add experience</button>
        </div>
        <AddExperience 
          :searchEnabled="searchEnabled"
          :addingReady="addingExperience"
          @handleSearch="handleSearch()"
          @disable="handleAddExperience()"
        />
      </div>

      <div class="bio-container">
        <p class="bio-text">{{bio}}</p>
      </div>
    </section>

    <div class="profile-edit">
      <router-link :to="`/account/${account}/${user}/settings`" class="edit-router-link">settings</router-link>
      <router-link :to="`/account/${account}/${user}/edit`" class="edit-router-link">edit profile</router-link>
    </div>
  </section>
</template>

<script>
import { GET_ACCOUNT } from '../queries'
import { Icon } from '@iconify/vue'
import AddExperience from '../sections/AddExperience.vue';

export default {
  name: 'Profile',

  components: {
    Icon,
    AddExperience
  },

  data() {
    return {
      name: '',
      email: '',
      imageFilePath: '',
      bio: 'no bio actually',
      workspaces: [],
      searchEnabled: false,
      addingExperience: false,
      addingOptionsUnfold: [false, false, false],
    }
  },

  computed: {
    user() {
      return localStorage.getItem('username')
    },

    account() {
      return localStorage.getItem('account')
    }
  },

  created() {
    this.fetcProfiledata()
  },

  methods: {
    async fetcProfiledata() {
      const id = localStorage.getItem('initialUser')
      const { data } = await this.$apollo.query({
        query: GET_ACCOUNT,
        variables: {userID: id}
      })

      this.name = data.getUser.profile.contact.name
      this.email = data.getUser.profile.contact.email
      this.imageFilePath = data.getUser.profile.contact.imageFilePath

      if (data.getUser.profile.bio !== "") {
        this.bio = data.getUser.profile.bio
      }
      
      if (data.getUser.experiences.length > 0) {
        this.experiences.push(...data.getUser.workspaces)
      }
    },

    handleAddExperience() {
      this.addingExperience = !this.addingExperience

      if (this.addingExperience) {
        setTimeout(() => {
          this.addingOptionsUnfold.forEach((_, index) => {
            setTimeout(() => {
              this.addingOptionsUnfold[index] = true
            }, index * 75)
          })
        }, 375);
      } else {
        this.addingOptionsUnfold.forEach((_, index) => {
          this.addingOptionsUnfold((_, index) => {
            this.addingOptionsUnfold[index] = false
          })
        })
      }
    },

    handleSearch() {
      this.searchEnabled = !this.searchEnabled
    },
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
  justify-content: space-between;
  height: 200px;
}

.profile-image-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: left;
  width: 20%;
  padding-left: 10rem;
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

.person-details-container {
  width: inherit;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.expierence-container {
  width: 30%;
  display: flex;
  flex-direction: column;
  align-items: left;
  min-height: 350px;
  padding-left: 10rem;
}

.experience {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  margin: 1rem;
}

.experience-description {
  font-size: 1.2rem;
  color: rgb(170, 225, 205);
}

.experience-star-rating-container {
  width: 1.2rem;
  height: 1.2rem;
  margin: 0.2rem;
  color: rgb(100, 103, 125);
}

.experience-star-rating-container.filled {
  color: rgb(240, 240, 100);
}

.add-experience-button-container {
  margin-top: auto;
  display: flex;
  justify-content: left;
}

.add-experience-button {
  background: inherit;
  border: none;
  color: rgb(100, 103, 125);
  font-size: 1rem;
  transition: all 0.3s ease;
  cursor: pointer
}

.add-experience-button:hover {
  color: rgb(170, 225, 205);
}

.bio-container {
  width: 50%;
  min-height: 350px;
  height: inherit;
  border: 1px solid rgb(170, 225, 205);
  border-radius: 5px;
  display: flex;
  justify-content: start;
  align-items: left;
}

.bio-text {
  margin: 1.2rem;
  font-size: 1.15rem;
  color: rgb(170, 225, 205);
}
</style>

