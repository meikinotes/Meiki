<script lang="ts">
    export let show: boolean = true
</script>

<span class="flex flex-row flex-grow text-toolbar-content justify-center font-light items-center">
    {#if show}
        <slot />
    {/if}
</span>
