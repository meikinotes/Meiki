<script lang="ts">
    import { createEventDispatcher } from "svelte"
    import { getUsername } from "$lib/api/user"
    import Toolbar from "$cmp/toolbar/Toolbar.svelte"
    import ToolbarButton from "$cmp/toolbar/Button.svelte"
    import ToolbarTitle from "$cmp/toolbar/Title.svelte"
    import ToolbarGroup from "$cmp/toolbar/Group.svelte"
    import SidebarIcon from "$cmp/icons/SidebarIcon.svelte"
    import EditIcon from "$cmp/icons/EditIcon.svelte"
    import PreviewIcon from "$cmp/icons/PreviewIcon.svelte"
    import UserIcon from "$cmp/icons/UserIcon.svelte"
    import ToolbarSyncIndicator from "$cmp/app/toolbar/SyncIndicator.svelte"
    import EditableLabel from "$cmp/EditableLabel.svelte"
    import ThemeSwitcher from "$cmp/toolbar/ThemeSwitcher.svelte"

    export let title: string

    export let explorerActive: boolean
    export let editorActive: boolean
    export let rendererActive: boolean

    export let showNoteActions: boolean
    export let changesNotSaved: boolean
    export let toolbarError: string

    const username = getUsername()
    const dispatchEvent = createEventDispatcher()
</script>

<Toolbar>
    <ToolbarButton isButtonChecked={explorerActive} name="sidebar" on:sidebar>
        <SidebarIcon />
    </ToolbarButton>

    <ToolbarGroup show={showNoteActions}>
        <span class="px-2" />
        <ToolbarButton isButtonChecked={editorActive} name="edit" on:edit>
            <EditIcon />
        </ToolbarButton>
        <ToolbarButton isButtonChecked={rendererActive} name="render" on:render>
            <PreviewIcon />
        </ToolbarButton>
    </ToolbarGroup>

    <ToolbarTitle show={showNoteActions}>
        <div data-cy="noteTitle" class="flex flex-row">
            <EditableLabel
                label={title}
                on:submit={(event) => {
                    dispatchEvent("rename", { newTitle: event.detail })
                }}
            />
            <ToolbarSyncIndicator
                show={showNoteActions || !!toolbarError}
                {changesNotSaved}
                error={toolbarError}
            />
        </div>
    </ToolbarTitle>
    <ThemeSwitcher />
    <ToolbarButton name="profile" label={username} on:profile>
        <UserIcon />
    </ToolbarButton>
</Toolbar>
