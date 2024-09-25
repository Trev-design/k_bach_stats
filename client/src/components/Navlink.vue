<template>
  <router-link class="navlink" :to="`/${dest}`">{{ dest }}</router-link>
</template>

<script>
export default {
  name: 'Navlink',
  props: {
    dest: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.navlink {
  position: relative;
  margin-right: 1.8rem;
  text-decoration: none;
  color: rgb(180, 200, 225);
  font-size: 1.4rem;
  font-weight: 500;
}

.navlink::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  background: rgb(180, 200, 225);
  bottom: -4px;
  left: 0;
  right: 0;
  transform: scale(0, 1);
  transition: transform 0.3s ease;
}

.navlink:hover::after {
  transform: scale(1, 1);
}
</style>