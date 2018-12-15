<template>
    <div class="container d-flex flex-grow-1">
        <chat/>
    </div>
</template>

<script>
import Chat from "@/components/Chat.vue";
import axios from "axios";

export default {
    name: "home",
    components: {
        Chat
    },
    data() {
        return {};
    },
    methods: {},
    computed: {
        token() {
            return this.$store.state.auth.token.slice(7).trim();
        }
    },
    mounted() {
        return this.$http
            .get("/api/health", {
                headers: {
                    Authorization: `Bearer ${this.token}`
                }
            })
            .then(() => {
                this.$connect(`ws://localhost:5050/api/ws?token=${this.token}`);
            });
    }
};
</script>
<style scoped lang="stylus">
.container
    margin-top 10px
</style>
