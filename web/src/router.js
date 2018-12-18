import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import Auth from "./components/Auth.vue";
import store from "./store";

Vue.use(Router);

let router = new Router({
    mode: "history",
    base: process.env.BASE_URL,
    routes: [
        {
            path: "/login",
            name: "login",
            component: Auth
        },
        {
            path: "/",
            name: "home",
            component: Home,
            meta: {
                protected: true
            }
        },
        {
            path: "/about",
            name: "about",
            // route level code-splitting
            // this generates a separate chunk (about.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: () =>
                import(/* webpackChunkName: "about" */ "./views/About.vue")
        },
        {
            path: "*"
        }
    ]
});

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.protected)) {
        if (store.getters.LOGGED_IN) {
            next();
            return;
        }
        next("/login");
    } else {
        next();
    }
});

export default router;
