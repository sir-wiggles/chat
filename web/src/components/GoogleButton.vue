<template>
    <div>
        <img
            class="container"
            v-if="!isLoading"
            :src="buttonState"
            @mouseenter="mouseenter"
            @mouseout="mouseout"
            @click="manualLogin"
        />
        <div v-else class="lds-ellipsis"><div></div><div></div><div></div><div></div></div>
    </div>
</template>
<script>
import { mapActions } from "vuex";

export default {
    name: "GoogleLogin",
    data() {
        return {
            config: {
                clientId:
                    "112260909295-mtovbmtnoc2r3gqhsipmjfcof8ik7uvj.apps.googleusercontent.com",
                scope: "https://www.googleapis.com/auth/userinfo.profile",
                response_type: "code"
            },
            loading: true,
            active: require("../../public/images/btn_google_signin_dark_normal_web.png"),
            disabled: require("../../public/images/btn_google_signin_dark_disabled_web.png"),
            focus: require("../../public/images/btn_google_signin_dark_focus_web.png"),
            pressed: require("../../public/images/btn_google_signin_dark_pressed_web.png"),
            state: "active"
        };
    },
    props: {},
    computed: {
        buttonState() {
            return this.$data[this.$data.state];
        },
        isPressed() {
            return this.$data.state === "pressed";
        },
        isDisabled() {
            return this.$data.state === "disabled";
        },
        isLoading() {
            return this.$data.loading;
        }
    },
    methods: {
        ...mapActions(["OAUTH2_GET_TOKEN"]),

        mouseenter() {
            if (this.isDisabled) {
                return;
            }
            if (!this.isPressed) {
                this.$data.state = "focus";
            }
        },

        mouseout() {
            if (this.isDisabled) {
                return;
            }
            if (!this.isPressed) {
                this.$data.state = "active";
            }
        },

        click() {
            if (this.isDisabled) {
                return;
            }
            this.$data.state = "pressed";
            this.login();
        },

        authorize(gapi, prompt) {
            this.$log.debug("authorizing");
            return new Promise((resolve, reject) => {
                this.$data.config.prompt = prompt;
                gapi.authorize(this.$data.config, function(response) {
                    if (response.error) {
                        reject(response.error);
                    }
                    resolve(response);
                });
            });
        },

        waitForGapi() {
            this.$log.debug("waiting for gapi");
            return new Promise(resolve => {
                setTimeout(() => {
                    this.$log.debug("...");
                    if (window.gapi.auth2) {
                        resolve(window.gapi.auth2);
                    } else {
                        resolve(this.waitForGapi());
                    }
                }, 250);
            });
        },

        autoLogin() {
            this.$log.debug("attempting autoLogin");
            this.$data.loading = true;
            this.waitForGapi()
                .then(gapi => {
                    this.$log.debug("gapi initialized");
                    return this.authorize(gapi, "none");
                })
                .then(({ code }) => {
                    this.$log.debug("autoLogin exchanging code for token");
                    return this.OAUTH2_GET_TOKEN(code);
                })
                .then(() => {
                    this.$router.push({ name: "home" });
                })
                .catch(error => {
                    this.$log.warn(`autoLogin ${error}`);
                })
                .finally(() => {
                    this.$data.loading = false;
                    this.$data.state = "active";
                });
        },

        manualLogin() {
            this.$log.debug("attempting manualLogin");
            this.$data.loading = true;
            this.waitForGapi()
                .then(gapi => {
                    return this.authorize(gapi);
                })
                .then(({ code }) => {
                    this.$log.debug("manualLogin exchanging code for token");
                    return this.OAUTH2_GET_TOKEN(code);
                })
                .then(() => {
                    this.$router.push({ name: "home" });
                })
                .catch(error => {
                    this.$log.warn(`manualLogin ${error}`);
                })
                .finally(() => {
                    this.$data.loading = false;
                    this.$data.state = "active";
                });
        }
    },

    beforeMount() {
        this.$store.dispatch("CHECK_TOKEN").then(valid => {
            if (!valid) {
                return this.autoLogin();
            } else {
                this.$router.push({ name: "home" });
            }
        });
    }
};
</script>
<style lang=stylus scoped>

.container
    cursor: pointer;

.lds-ellipsis
  display inline-block
  position relative
  width 64px
  height 64px
  top -15px

.lds-ellipsis div
  position absolute
  top 27px
  width 11px
  height 11px
  border-radius 50%
  background #0ff
  animation-timing-function cubic-bezier(0, 1, 1, 0)



.lds-ellipsis div:nth-child(1)
    left 6px
    animation lds-ellipsis1 0.6s infinite
    background #ea4335 // red

.lds-ellipsis div:nth-child(2)
    left 6px
    animation lds-ellipsis2 0.6s infinite
    background #fbbc2b // yellow

.lds-ellipsis div:nth-child(3)
    left 26px
    animation lds-ellipsis2 0.6s infinite
    background #33a853 // green

.lds-ellipsis div:nth-child(4)
    left 45px
    animation lds-ellipsis3 0.6s infinite
    background #4285f4 // blue

@keyframes lds-ellipsis1
  0%
    transform scale(0)

  100%
    transform scale(1)


@keyframes lds-ellipsis3
  0%
    transform scale(1)

  100%
    transform scale(0)


@keyframes lds-ellipsis2
  0%
    transform translate(0, 0)

  100%
    transform translate(19px, 0)

</style>
