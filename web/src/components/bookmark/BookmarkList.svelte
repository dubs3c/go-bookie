<script>
  import PureBookmarkList from './PureBookmarkList.svelte';
  import { bookmarkStore } from '../../store';
  function onPinTask(event) {
    taskStore.pinTask(event.detail.id);
  }
  function onArchiveTask(event) {
    taskStore.archiveTask(event.detail.id);
  }
</script>

<div>
  <PureBookmarkList
    bookmarks={$bookmarkStore}
    on:onPinTask={onPinTask}
    on:onArchiveTask={onArchiveTask}
  />
</div>