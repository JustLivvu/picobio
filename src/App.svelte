<script lang="ts">
    import { createRouter, View } from "@shaun/svelte-router";
    import Home from "./views/Home.svelte";
    import Panel from "./views/Panel.svelte";

    createRouter({
        mode: "web",
        routes: [
            { path: "/", component: Home },
            { path: "/panel", component: Panel },
            { path: "*", component: Home }
        ]
    });
</script>

<View />