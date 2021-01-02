<template>
    <v-card
        class="mx-auto"
    >
      <v-card-title style="background:  #74767a" >
        <v-icon color="white" >mdi-server-network</v-icon>
        <span class="headline white--text">KubeMQ Servers</span>
        <v-spacer></v-spacer>
      </v-card-title>
      <v-row mt-0 mb-0>
      <v-container class="lighten-5 ">
        <v-row mb-0>
        <v-col>
          <v-card-title >
            <span class="headline black--text">Incoming Traffic</span>
            <v-spacer></v-spacer>
          </v-card-title>
          <v-container class="lighten-5 ">
            <v-row class="mb-6" >
              <v-col v-for="item in items.total_in" :key="item.key">
                <DataCard v-bind:title="item.key" v-bind:value="item.value"></DataCard>
              </v-col>
            </v-row>

          </v-container>
        </v-col>
        <v-col>
          <v-card-title >
            <span class="headline black--text">Outgoing Traffic</span>
            <v-spacer></v-spacer>
          </v-card-title>
          <v-container class="lighten-5 ">
            <v-row class="mb-6" >
              <v-col v-for="item in items.total_out" :key="item.key">
                <DataCard v-bind:title="item.key" v-bind:value="item.value"></DataCard>
              </v-col>
            </v-row>
          </v-container>
        </v-col>
        </v-row>
        <v-row >

        </v-row>
        </v-container>
      </v-row>
      <v-row mt-0>
        <v-container class="lighten-5 ">
          <v-data-table
              :headers="headers"
              :items="items.list"
              :hide-default-footer="true"
              class="elevation-0"
          ></v-data-table>

        </v-container>

      </v-row>
    </v-card>

</template>

<script>
import DataCard from "@/components/DataCard";
import axios from "axios";
import EventBus from "@/event-bus";

export default {
  name: "Kubemq",
  components: {DataCard},
  props:[
    'baseUrl','pollInterval'
  ],
  data () {
    return {
      items: {},
      headers: [
        {
          text: 'Source',
          value: 'title',
        },
        {text: 'Channels', value: 'channels'},
        {text: 'Messages', value: 'messages'},
        {text: 'Volume', value: 'volume'},
        {text: 'Pending', value: 'pending'},
        {text: 'Errors', value: 'errors'},
        {text: 'CPU (Cores)', value: 'cpu'},
        {text: 'CPU (%)', value: 'cpu_utilization'},
        {text: 'Memory', value: 'memory'},
        {text: 'Memory (%)', value: 'memory_utilization'},
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
    pollData () {
      this.polling = setInterval(() => {
        this.getData()
      }, this.pollInterval)
    },
    getData: function () {
      axios
          .get(this.baseUrl+`/kubemq`)
          .then(response =>  this.items=response.data)
          .catch(()=> this.items={})
    },
    clear: function (){
      this.items= {}
    }

  },
}
</script>
