<template>
  <v-container>
    <h5 class="grey--text">{{ title }}</h5>
    <h3 class="grey--text text--darken-2">{{ value }}</h3>
  </v-container>
  </template>

<script>
export default {
  name: "DataCard",
  props:[
      'title','value'
  ],

}
</script>

<style scoped>

</style>
