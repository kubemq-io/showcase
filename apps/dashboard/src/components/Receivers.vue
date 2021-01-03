<template>
  <v-card
      class="mx-auto"
  >
    <v-card-title >
      <v-icon color="white"  >mdi-arrow-down</v-icon>
      <span class="headline white--text">Receivers</span>
      <v-spacer></v-spacer>
    </v-card-title>
    <v-container class="lighten-5 ">
      <v-row
          class="mb-0"
      >
        <v-col v-for="item in items.total" :key="item.key">
          <DataCard v-bind:title="item.key" v-bind:value="item.value"></DataCard>
        </v-col>
      </v-row>
    </v-container>
    <v-data-table
        :headers="headers"
        :items="items.list"
        :hide-default-footer="true"
        class="elevation-0"
    ></v-data-table>
  </v-card>
</template>

<script>

import DataCard from "@/components/DataCard";
import axios from "axios";
import EventBus from "@/event-bus";
export default {
  name: "Receivers",
  props:[

  ],
  components: {DataCard},
  data: function () {
    return {
      items: {},
      headers: [
        {
          text: 'Source',
          value: 'title',
        },
        { text: 'Clients', value: 'clients' },
        { text: 'Messages', value: 'messages' },
        { text: 'Volume', value: 'volume' },
        { text: 'Errors', value: 'errors' },
      ],
      polling: null
    }
  },
  mounted() {
    EventBus.$on('clear',this.clear )
  },
  created() {
    this.getData();
    this.pollData();
  },

  methods: {
    pollData() {
      this.polling = setInterval(() => {
        this.getData()
      }, this.POLL_INTERVAL)
    },
    getData: function () {
      axios
          .get(this.API_SERVER_URL + `/receivers`)
          .then(response => this.items = response.data)
          .catch(() => this.items = {})
    },
    clear: function (){
      this.items= {}
    }

  },
}
</script>

<style scoped>
.v-card__title {
  background: #FFB064
}

</style>
