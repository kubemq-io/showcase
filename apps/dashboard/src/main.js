import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import axios from 'axios'
import VueAxios from 'vue-axios'
import '@mdi/font/css/materialdesignicons.css'

Vue.config.productionTip = false
export const bus = new Vue();
new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')

Vue.use(VueAxios, axios)
