<script>
    import Bookmark from './Bookmark.svelte';
    export let loading = false;
    export let bookmarks = [];
  
    // reactive declarations (computed prop in other frameworks)
    $: noBookmarks = bookmarks.length === 0;
    $: emptyBookmarks = bookmarks.length === 0 && !loading;
  </script>
  {#if loading}
    <div class="list-items">loading...</div>
  {/if}
  {#if emptyBookmarks}
    <h3>You don't have any bookmarks yet!</h3>
  {/if}
  {#each bookmarks as bookmark}
    <Bookmark {bookmark}/>
  {/each}