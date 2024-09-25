<template>
  <div class="sectionlink-group-container">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'SectionlinkGroup',
}
</script>

<style>
.sectionlink-group-container {
  display: flex;
  flex-direction: row;
  align-items: end;
  justify-content: center;
  margin-right: 0 2rem;
}
</style>