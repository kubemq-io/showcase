<template>
  <v-app>
    <v-app-bar
      app
      color="#4C809E"
      dark
    >
     <h1>Showcase Dashboard</h1>

      <v-spacer></v-spacer>

      <v-btn
        href="https://kubemq.io"
        target="_blank"
        text
      >
        <span class="mr-2">KubeMQ.io</span>

      </v-btn>
    </v-app-bar>

    <v-main>
      <Main></Main>
    </v-main>
    <v-footer
        dark
        padless
    >
      <v-card
          class="flex"
          flat
          tile
      >


        <v-card-text class="py-2 white--text text-center">
          {{ new Date().getFullYear() }} — <strong>KubeMQ</strong>
        </v-card-text>
      </v-card>
    </v-footer>
  </v-app>
</template>

<script>

import Main from "@/pages/Main";

export default {
  name: 'App',

  components: {
    Main,
  },

  created() {
    document.title='KubeMQ Showcase Dashboard'
},
  data: () => ({
    //
  }),
};
</script>
