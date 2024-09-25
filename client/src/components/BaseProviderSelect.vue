<template>
  <div class="provider-select-container">
    <div class="provider-select-button">
      <div class="provider-select-button-text-container">
        <Icon 
          icon="mdi:google"
          class="provider-select-button-icon"
        />
        <p class="provider-select-button-text">Signup with google</p>
      </div>
    </div>
    <div class="provider-select-button">
      <div class="provider-select-button-text-container">
        <Icon 
          icon="mdi:facebook"
          class="provider-select-button-icon"
        />
        <p class="provider-select-button-text">Signup with facebook</p>
      </div>
    </div>
    <div class="provider-select-button">
      <div class="provider-select-button-text-container">
        <Icon 
          icon="mdi:github"
          class="provider-select-button-icon"
        />
        <p class="provider-select-button-text">Signup with github</p>
      </div>
    </div>
    <div class="provider-select-button">
      <div class="provider-select-button-text-container">
        <Icon 
          icon="mdi:gitlab"
          class="provider-select-button-icon"
        />
        <p class="provider-select-button-text">Signup with gitlab</p>
      </div>
    </div>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'
export default {
  name: 'BaseProviderSelect',

  components: {
    Icon
  }
}
</script>

<style scoped>
.provider-select-container {
  width: 72%;
  display: flex;
  flex-direction: column;
  justify-content: left;
  align-items: left;
}

.provider-select-button {
  width: 100%;
  cursor: pointer;
  margin: 0.2rem 0;
  padding: 0.1rem 0.5rem;
  border-radius: 3px;
}

.provider-select-button-text-container {
  display: flex;
  justify-content: left;
  align-items: center;
}

.provider-select-button-text {
  padding: 0 1rem;
  font-size: 1.15rem;
  color: #03fc77;
}

.provider-select-button-icon {
  width: 1.2rem;
  height: 1.2rem;
  color: #03fc77;
}

.provider-select-button-text-container:hover > .provider-select-button-text {
  color: aquamarine;
  transition: .3s ease;
}

.provider-select-button-text-container:hover > .provider-select-button-icon {
  color: aquamarine;
  width: 1.4rem;
  height: 1.4rem;
  transition: .3s ease;
}
</style>