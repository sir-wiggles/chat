import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import VueNativeSock from "vue-native-websocket";

import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "./stylesheet.css";

Vue.config.productionTip = false;

Vue.use(BootstrapVue);

Vue.use(VueNativeSock, "ws://localhost:5050/ws", {
  store,
  connectManually: true,
  format: "json"
});

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
