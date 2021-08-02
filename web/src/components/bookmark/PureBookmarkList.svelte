<script>
    import Bookmark from './Bookmark.svelte';
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    export let loading = false;
    export let bookmarks = [];
  
    function forward(event) {
		  dispatch('ArchiveTask', event.detail);
	  }

    function forwardDelete(event) {
		  dispatch('onDeleteBookmark', event.detail);
	  }

    // reactive declarations (computed prop in other frameworks)
    $: emptyBookmarks = bookmarks.length === 0 && !loading;
  </script>
  {#if loading}
    <div class="list-items">loading...</div>
  {/if}
  {#if emptyBookmarks}
    <h3>You don't have any bookmarks yet!</h3>
  {/if}
  {#each bookmarks as bookmark}
    <Bookmark 
    on:onDeleteBookmark={forwardDelete}
    on:ArchiveTask={forward} {bookmark}/>
  {/each}