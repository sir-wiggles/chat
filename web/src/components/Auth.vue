<template>
    <div class="container">
        <google-button />
    </div>
</template>

<script>
import GoogleButton from "@/components/GoogleButton";

export default {
    name: "login",
    components: {
        GoogleButton
    }
};
</script>
<style lang="stylus" scoped>
.container
    height 100vh
    display flex
    align-items center
    justify-content center
</style>
