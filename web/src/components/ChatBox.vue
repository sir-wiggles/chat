<template>
  <div class="container d-flex flex-column">
    <registration />
    <message-list
      class="message-list flex-grow-1"
      :messages="messages"
    ></message-list>
    <b-form-input
      class="input"
      v-model="input"
      @keyup.enter.native="send"
    ></b-form-input>
  </div>
</template>

<script>
import { mapState } from "vuex";
import MessageList from "@/components/MessageList.vue";
import Registration from "@/components/Registration.vue";

export default {
  name: "chat-box",
  components: {
    MessageList,
    Registration
  },
  data() {
    return {
      input: ""
    };
  },
  methods: {
    send() {
      if (this.input.trim().length === 0) {
        return;
      }
      this.$socket.send(this.input);
      this.input = "";
    },
    close() {
      this.$refs.registration.hide();
    }
  },
  computed: {
    ...mapState({
      messages: state => state.messages
    })
  },
  mounted() {
    //this.$connect("ws://localhost:5050/ws");
  },
  destroyed() {
    this.$disconnect();
  }
};
</script>

<style scoped lang="stylus">

.message-list
  overflow scroll
  overflow-x hidden
  overflow-y auto

  color: #495057;
  background-color: #fff;
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
  margin 10px 0 0 0

.input
  margin 10px 0
  min-height 38px
</style>
