<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<script lang="ts">
import "./bookmark.css"
import { createEventDispatcher } from 'svelte';
import { DeleteBookmark, ToggleStatus } from "./../../actions/BookmarkAction.svelte"
import type {Bookmark} from "./../../types/Bookmark"
import Tag from "../tag/Tag.svelte";

const dispatch = createEventDispatcher();

// event handler for Pin Task
async function HandleDeleteBookmark(BookmarkID) {
    await DeleteBookmark(BookmarkID)
    dispatch('onDeleteBookmark', {
        id: BookmarkID,
    });
}

async function ChangeStatus(BookmarkID: number, deleted: boolean, archived: boolean) {
    await ToggleStatus(BookmarkID, deleted, archived)
    dispatch('onDeleteBookmark', {
        id: BookmarkID,
    });
}

// Bookmark props
export let bookmark: Bookmark

</script>

<div class="bookmark">
    <h4><a href="{bookmark.url}" target="_blank" rel="noopener noreferrer">{bookmark.title}</a>
    </h4>
    <small>{bookmark.url}</small>
    <div class="actions">
        <a href="/bookmark/{bookmark.id}" class="view">
            <i class="fas fa-eye"></i> View
        </a>
        {#if bookmark.deleted}
            <button on:click={() => HandleDeleteBookmark(bookmark.id)}  class="delete {bookmark.deleted == true ? "red": ""}">
                <i class="fas fa-trash-alt"></i> Delete
            </button>

            <button on:click={() => ChangeStatus(bookmark.id, false, bookmark.archived)}  class="delete blue">
                <i class="fas fa-trash-alt"></i> Restore
            </button>
        {:else}
            <button on:click={() => ChangeStatus(bookmark.id, true, bookmark.archived)}  class="delete {bookmark.deleted == true ? "red": ""}">
                <i class="fas fa-trash-alt"></i> Trash
            </button>
        {/if}


        {#if bookmark.deleted != true }
            <button on:click={() => ChangeStatus(bookmark.id, bookmark.deleted, !bookmark.archived)} class="archive {bookmark.archived == true ? "green": ""}">
                {#if bookmark.archived}
                    <i class="fas fa-bookmark"></i> Un-archive
                {:else}
                    <i class="fas fa-bookmark"></i> archive
                {/if}
            </button>
        {/if}

        {#if bookmark.tags != ""}
        <br />
            {#each bookmark.tags.split(",") as tag}
                <Tag value="{tag}"/>
            {/each}
        {/if}

    </div>
</div>
<div class="border"></div>