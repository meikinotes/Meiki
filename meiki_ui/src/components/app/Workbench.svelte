<script context="module" lang="ts">
    function getEditorClass(editorActive: boolean, rendererActive: boolean) {
        if (editorActive && rendererActive) return "flex-1"
        if (editorActive && !rendererActive) return "w-full"

        return "hidden"
    }

    function getRendererClass(editorActive: boolean, rendererActive: boolean) {
        if (editorActive && rendererActive) return "flex-1 w-1/2"
        if (rendererActive && !editorActive) return "max-w-[800px]"

        return "hidden"
    }
</script>

<script lang="ts">
    import { afterUpdate, createEventDispatcher } from "svelte"
    import Editor from "$cmp/Editor.svelte"
    import Renderer from "$cmp/Renderer.svelte"
    import Seperator from "$cmp/Seperator.svelte"

    export let fontSize = 16
    export let showWorkbench: boolean
    export let editorActive: boolean
    export let rendererActive: boolean

    let text = ""
    let editor: Editor
    const dispatchEvent = createEventDispatcher()

    function onEditorChange() {
        text = editor.getValue()
        dispatchEvent("textChange", { text })
    }

    function focus() {
        if (editor) editor.focus()
    }

    export function setText(newText: string) {
        text = newText
        if (editor) {
            editor.setValue(newText)
        }
    }

    export function getText() {
        return editor.getValue()
    }

    afterUpdate(focus)
</script>

<div class="flex flex-grow justify-center items-center bg-base-0 overflow-x-scroll">
    <div class="flex justify-center h-full w-full">
        {#if showWorkbench}
            <div class={getEditorClass(editorActive, rendererActive)}>
                <Editor
                    bind:this={editor}
                    {fontSize}
                    onChange={onEditorChange}
                    editorId="workbenchEditor"
                    initialText={text}
                />
            </div>
            {#if editorActive && rendererActive}
                <Seperator />
            {/if}
            <div class={getRendererClass(editorActive, rendererActive)}>
                <Renderer {text} />
            </div>
        {/if}
    </div>
</div>
