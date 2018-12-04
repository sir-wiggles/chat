<template>
  <b-modal ref="registration" hide-footer>
    <div class="group d-flex">
      <label class="label" :for="username">Username</label>
      <b-form-input
        :id="username"
        v-model="username"
        :state="usernameState"
      ></b-form-input>
    </div>
    <div class="group d-flex">
      <label class="label" :for="email">Email</label>
      <b-form-input :id="email" v-model="email"></b-form-input>
    </div>
    <div class="group d-flex">
      <label class="label" :for="avatar">Avatar</label>
      <b-form-input :id="avatar" v-model="avatar"></b-form-input>
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
      avatar: "",
      email: "",
      username: ""
    };
  },
  methods: {
    register() {
      if (this.usernameState) {
        let payload = {
          username: this.username,
          email: this.email,
          avatar: this.avatar
        };
        this.$store.dispatch("REGISTER_USER", payload).then(resp => {
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
    this.$refs.registration.show();
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
