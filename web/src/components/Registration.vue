<template>
    <b-modal ref="registration" hide-footer>
        <div class="group d-flex">
            <label class="label" :for="username">Username</label>
            <b-form-input
                :id="username"
                v-model="username"
                :state="usernameState"
                />
        </div>
        <div class="group d-flex">
            <label class="label" :for="password">Password</label>
            <b-form-input :id="password" v-model="password" type="password"/>
        </div>
        <b-btn class="mt-3" @click="register">Go</b-btn>
    </b-modal>
</template>

<script>
export default {
    name: "registrationModal",
    components: {},
    data() {
        return {
            password: "",
            username: ""
        };
    },
    methods: {
        register() {
            if (this.usernameState) {
                let payload = {
                    username: this.username,
                    password: this.password
                };
                this.$store
                    .dispatch("AUTHENTICATE_USER", payload)
                    .then(resp => {
                        if (resp === true) {
                            this.$refs.registration.hide();
                        }
                    });
            }
        }
    },
    computed: {
        usernameState() {
            if (this.username.trim().length >= 2) {
                return true;
            }
            return false;
        }
    },
    mounted() {
        if (localStorage.username && localStorage.username.length > 0) {
            this.$store.dispatch("AUTHENTICATE_USER", {
                username: localStorage.username,
                password: localStorage.password
            });
        } else {
            this.$refs.registration.show();
        }
    }
};
</script>

<style scoped lang="stylus">
.label
    width 100px
    text-align left
    padding-top 6px

.group
    margin 3px 0 3px 0
</style>
