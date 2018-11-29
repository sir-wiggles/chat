<template>
  <b-list-group-item flush href="#">
    <div v-if="isAuthorsMessage" class="d-flex flex-row-reverse">
      <img class="avatar" :src="avatar" />
      <div class="flex-column message-area message-right">
        <div class="d-flex flex-row justify-content-end">
          <div class="author ">{{ name }}</div>
          <div class="time">{{ time }}</div>
        </div>
        <div class="author-text">
          <pre>{{ text }}</pre>
        </div>
      </div>
    </div>
    <div
      v-else-if="isSystemMessage"
      class="d-flex flex-row justify-content-center"
    >
      <img class="avatar" :src="avatar" />
      <div class="flex-column message-area message-center">
        <div class="d-flex flex-row">
          <div class="author">Admin</div>
          <div class="time">{{ time }}</div>
        </div>
        <div class="text">
          <pre>{{ text }}</pre>
        </div>
      </div>
    </div>
    <div v-else class="d-flex flex-row">
      <img class="avatar" :src="avatar" />
      <div class="flex-column message-area message-left">
        <div class="d-flex flex-row">
          <div class="author">{{ name }}</div>
          <div class="time">{{ time }}</div>
        </div>
        <div class="text">
          <pre>{{ text }}</pre>
        </div>
      </div>
    </div>
  </b-list-group-item>
</template>

<script>
import moment from "moment";
import bListGroupItem from "bootstrap-vue/es/components/list-group/list-group-item";

export default {
  components: {
    bListGroupItem
  },
  data() {
    return {};
  },
  props: {
    message: {
      type: Object,
      required: true
    }
  },
  methods: {},
  computed: {
    text() {
      return this.message.text.join("\n");
    },
    name() {
      return this.message.author.name;
    },
    id() {
      return this.$store.state.author.id;
    },
    time() {
      return moment(this.message.time).format("HH:mm");
    },
    avatar() {
      return this.message.author.avatar;
    },
    isSystemMessage() {
      return { system: true, initialize: true }[this.message.type];
    },
    isAuthorsMessage() {
      return (
        this.message.type === "message" && this.message.author.id === this.id
      );
    }
  }
};
</script>
<style scoped lang="stylus">
.list-group-item
  padding 0
  border 0

.avatar
  margin 6px
  width 36px
  height 36px

.author
  font-weight bold
  padding 0 10px 0 2px

.time
  font-size 11px
  padding 1px 0 0 0

pre
  margin 0
  white-space pre-line

.text
  text-align left

.author-text
  text-align right

.message-area
  border-radius 6px
  padding 4px

.message-left
  background lightgreen
  margin-right 48px

.message-center
  background lightblue

.message-right
  background lightgrey
  margin-left 48px
</style>
