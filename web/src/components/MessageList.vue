<template>
  <b-list-group class="container" flush ref="scrollList">
    <Message
      v-for="(message, idx) in messages"
      :message="message"
      :key="idx"
      class="message"
    />
  </b-list-group>
</template>
<script>
import bListGroup from "bootstrap-vue/es/components/list-group/list-group";
import Message from "./Message.vue";

export default {
  components: {
            Message,
    bListGroup
        },
        data() {
            return {
                previousMessage: null
    };
  },
        props: {
    messages: {
      type: Array,
                required: true
            }
  },
        methods: {
    scrollDown() {
                this.$refs.scrollList.scrollTop = this.$refs.scrollList.scrollHeight;
    }
  },
        updated() {
            this.scrollDown();
  },
        computed: {}
};
</script>

<style scoped lang="stylus">

.container
    border 1px solid grey
    border-radius: 3px
    box-shadow 1px 1px 11px grey

.message
    margin 3px 0

/* width */
::-webkit-scrollbar
    width 10px

/* Track */
::-webkit-scrollbar-track
    background #f1f1f1

/* Handle */
::-webkit-scrollbar-thumb
    background #888

/* Handle on hover */
::-webkit-scrollbar-thumb:hover
    background: #555
</style>
