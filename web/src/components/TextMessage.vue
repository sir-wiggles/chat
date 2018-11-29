<template>
  <div class="sc-message--text" :style="messageColors">
    <p v-html="messageText"></p>
    <p
      v-if="data.meta"
      class="sc-message--meta"
      :style="{ color: messageColors.color }"
    >
      {{ data.meta }}
    </p>
  </div>
</template>

<script>
export default {
  props: {
    data: {
      type: Object,
      required: true
    },
    messageColors: {
      type: Object,
      required: true
    },
    messageStyling: {
      type: Boolean
      //required: true
    }
  },
  computed: {
    messageText() {
      return this.data.text;
      //const escaped = escapeGoat.escape(this.data.text);
      //return Autolinker.link(this.messageStyling ? fmt(escaped) : escaped, {
      //  className: "chatLink",
      //  truncate: { length: 50, location: "smart" }
      //});
    }
  }
};
</script>

<style scoped>
a.chatLink {
  color: inherit !important;
}
</style>
