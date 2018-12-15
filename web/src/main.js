import "setimmediate";
import axios from "axios";
import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import VueNativeSock from "vue-native-websocket";
import { createLogger, format, transports } from "winston";
import "bootstrap/dist/css/bootstrap.css";
import "bootstrap-vue/dist/bootstrap-vue.css";
import "./stylesheet.css";

import App from "./App.vue";
import store from "./store";
import router from "./router";

Vue.config.productionTip = false;

const http = axios.create({
    baseURL: "http://localhost:5050",
    timeout: 1000
    //headers: {'X-Custom-Header': 'foobar'}
});

const { combine, timestamp, printf } = format;
const fmt = printf(info => {
    return `${info.timestamp} [${info.level}] ${info.message}`;
});

const logger = createLogger({
    level: "debug",
    format: combine(format.colorize(), timestamp(), fmt),
    transports: [new transports.Console()]
});

Vue.prototype.$http = http;
Vue.prototype.$log = logger;

store.$http = http;
store.$log = logger;

Vue.use(BootstrapVue);
Vue.use(VueNativeSock, `ws://localhost:8080/ws`, {
    store,
    connectManually: true,
    format: "json"
});

new Vue({
    router,
    store,
    render: h => h(App)
}).$mount("#app");
