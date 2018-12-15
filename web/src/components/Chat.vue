<template>
    <div class="container d-flex flex-grow-1">
        <chat-box />
    </div>
</template>
<script>
import ChatBox from "@/components/ChatBox.vue";
export default {
    name: "chat",
    components: {
        ChatBox
    },
    data() {
        return {};
    },
    methods: {}
};
</script>
<style>
</style>
