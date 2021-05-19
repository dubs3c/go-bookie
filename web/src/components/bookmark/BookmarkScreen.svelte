
<script>
    import BookmarkList from './BookmarkList.svelte';
    export let error = false;
</script>
  

{#if error}

    <div class="title-message">Oh no!</div>
    <div class="subtitle-message">Something went wrong</div>

{:else}

    <BookmarkList />

{/if}