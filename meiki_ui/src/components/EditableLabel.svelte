<script lang="ts">
    // Stolen and improved upon https://svelte.dev/repl/29c1026dda3c47a187bd21afa0782df1?version=3.48.0

    import { createEventDispatcher } from "svelte"

    export let label: string

    const dispatch = createEventDispatcher()
    let editing: boolean = false
    let value: string = label

    function edit() {
        value = label
        editing = true
    }

    function submit() {
        if (value != label && value.length != 0) {
            dispatch("submit", value)
            label = value
        }

        value = label
        editing = false
    }

    function keydown(event: KeyboardEvent) {
        if (event.key == "Escape") {
            event.preventDefault()
            value = label
            editing = false
        }
    }

    function focus(element: HTMLInputElement) {
        element.focus()
        element.select()
    }
</script>

{#if editing}
    <form on:submit|preventDefault={submit} on:keydown={keydown}>
        <input
            class="border-none bg-transparent text-inherit focus:outline-none focus:border-primary-focus focus:ring-primary-focus focus:ring-2 rounded-lg px-2"
            bind:value
            on:blur={submit}
            use:focus
        />
    </form>
{:else}
    <div on:click={edit}>
        {label}
    </div>
{/if}
