import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import VueAxios from 'vue-axios'
import '@mdi/font/css/materialdesignicons.css'


const getRuntimeConfig = async () => {
  const runtimeConfig = await fetch('/runtimeConfig.json');
  return await runtimeConfig.json()
}

getRuntimeConfig().then(function(json) {
  Vue.config.productionTip = false
  Vue.mixin({
   data() {
      return {
        API_SERVER_URL: json.API_SERVER_URL,
        POLL_INTERVAL: json.POLL_INTERVAL
      }
    },
  });

  Vue.config.productionTip = false
  Vue.use(VueAxios, axios)
  new Vue({
    vuetify,
    render: h => h(App)
  }).$mount('#app')

});