<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<script lang="ts">
import "./bookmark.css"
import { createEventDispatcher } from 'svelte';
import { DeleteBookmark } from "./../../actions/BookmarkAction.svelte"
import type {Bookmark} from "./../../types/Bookmark"
import Tag from "../tag/Tag.svelte";

const dispatch = createEventDispatcher();

// event handler for Pin Task
function HandleDeleteBookmark(BookmarkID) {
    DeleteBookmark(BookmarkID)
    dispatch('onDeleteBookmark', {
        id: BookmarkID,
    });
}

// event handler for Archive Task
function HandleArchiveBookmark(BookmarkID) {
    dispatch('ArchiveTask', {
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
        <a href="\#" class="view">
            <i class="fas fa-eye"></i> View
        </a>
        <button on:click={() => HandleDeleteBookmark(bookmark.id)}  class="delete {bookmark.deleted == true ? "red": ""}">
            <i class="fas fa-trash-alt"></i> {bookmark.deleted == true ? "Removed": "Trash"}
        </button>

        <button on:click={() => HandleArchiveBookmark(bookmark.id)} class="archive {bookmark.archived == true ? "green": ""}">
            <i class="fas fa-bookmark"></i> {bookmark.archived == true ? "Archived": "Archive"}
        </button>

        {#if bookmark.tags != ""}
        <br />
            {#each bookmark.tags.split(",") as tag}
                <Tag value="{tag}"/>
            {/each}
        {/if}

    </div>
</div>
<div class="border"></div>