<template>
  <v-card elevation="0" >
  <v-container class="white">
    <v-row
        class="mb-6"
    >
      <v-col>
        <Senders base-url="http://localhost:8085" poll-interval="5000"></Senders>
      </v-col>
      <v-col >
        <Receivers base-url="http://localhost:8085" poll-interval="5000"></Receivers>
      </v-col>
    </v-row>
    <v-row class="mb-6">
      <v-col
      >
        <Kubemq base-url="http://localhost:8085" poll-interval="5000"></Kubemq>
      </v-col>
    </v-row>
  <v-row>
    <v-col>
      <v-card class="white" elevation="0" >
      <v-card-actions >
        <v-spacer></v-spacer>
        <v-btn @click="clear()" style="background: #4C809E;color: white" >Clear</v-btn>
      </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
  </v-container>

  </v-card>
</template>

<script>
import Senders from "@/components/Senders";
import Receivers from "@/components/Receivers";
import Kubemq from "@/components/Kubemq";
import axios from "axios";
import EventBus from '../event-bus';

export default {
name: "Main",
  components: {Kubemq, Receivers, Senders},

  methods:{
  clear () {
    EventBus.$emit('clear', null);
    axios
        .post("http://localhost:8085"+`/clear`)



  }
  },
}
</script>

<style scoped>

</style>


