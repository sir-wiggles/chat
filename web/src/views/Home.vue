<template>
  <div class="container d-flex flex-grow-1"><chat-box /></div>
</template>

<script>
import ChatBox from "@/components/ChatBox.vue";

export default {
  name: "home",
        components: {
    ChatBox
        },
  data() {
            return {};
  },
        methods: {}
};
</script>
<style scoped lang="stylus"></style>
