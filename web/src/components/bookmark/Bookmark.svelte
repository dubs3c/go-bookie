<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.1/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
<script>
    import "./bookmark.css"
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();

    // event handler for Pin Task
    function DeleteBookmark(BookmarkID) {
    dispatch('onDeleteBookmark', {
        id: BookmarkID,
    });
    }

    // event handler for Archive Task
    function ArchiveBookmark(BookmarkID) {
        dispatch('ArchiveTask', {
            id: BookmarkID,
        });
    }

    // Bookmark props
    export let bookmark = {
        id: '',
        title: '',
        url: "",
        description: "",
        image: "",
        archived: false,
        deleted: false
    };

</script>

<div class="bookmark">
    <h4><a href="{bookmark.url}" target="_blank" rel="noopener noreferrer">{bookmark.title}</a>
    </h4>
    <small>{bookmark.url}</small>
    <div class="actions">
        <a href="\#" class="view">
            <i class="fas fa-eye"></i> View
        </a>
        <button on:click={DeleteBookmark(bookmark.id)}  class="delete {bookmark.deleted == true ? "red": ""}">
            <i class="fas fa-trash-alt"></i> {bookmark.deleted == true ? "Deleted": "Delete"}
        </button>

        <button on:click={ArchiveBookmark(bookmark.id)} class="archive {bookmark.archived == true ? "green": ""}">
            <i class="fas fa-bookmark"></i> {bookmark.archived == true ? "Archived": "Archive"}
        </button>

    </div>
</div>
<div class="border"></div>