<script lang="ts">
    import Home from "./views/Home.svelte";
    import Panel from "./views/Panel.svelte";
    import NotFound from "./views/404.svelte";

    let currentPath = $state(typeof window !== "undefined" ? window.location.pathname : "/");

    $effect(() => {
        const updatePath = () => {
            currentPath = window.location.pathname;
        };

        window.addEventListener("popstate", updatePath);
        window.addEventListener("hashchange", updatePath);

        return () => {
            window.removeEventListener("popstate", updatePath);
            window.removeEventListener("hashchange", updatePath);
        };
    });
</script>

{#if currentPath === "/" || (typeof window !== "undefined" && window.location.hash === "#/")}
    <Home />
{:else if currentPath === "/panel" || (typeof window !== "undefined" && window.location.hash === "#/panel")}
    <Panel />
{:else}
    <NotFound />
{/if}
