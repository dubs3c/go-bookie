<script>
  import PureBookmarkList from './PureBookmarkList.svelte';
  import { bookmarkStore } from '../../store';


  function onDeleteBookmark(event) {
    bookmarkStore.deleteBookmark(event.detail.id);
  }
  function onArchiveTask(event) {
    bookmarkStore.archiveTask(event.detail.id);
  }
</script>

<div>
  <PureBookmarkList
    bookmarks={$bookmarkStore}
    on:onDeleteBookmark={onDeleteBookmark}
    on:ArchiveTask={onArchiveTask}
  />
</div>