
import { writable } from 'svelte/store';

const Bookmarks = () => {
  // creates a new writable store populated with some initial data
  const { subscribe, update, set } = writable([]);

  return {
    subscribe,
    update,
    set,
    // method to archive a task, think of a action with redux or Vuex
    archiveBookmark: id =>
      update(bookmarks =>
        bookmarks.map(bookmark => (bookmark.id === id ? { ...bookmark, archived: !bookmark.archived } : bookmark))
      ),
    // method to archive a bookmark, think of a action with redux or Vuex
    deleteBookmark: id =>
      update(bookmarks =>
        bookmarks.map(bookmark => (bookmark.id === id ? { ...bookmark, deleted: !bookmark.deleted } : bookmark))
      ),
  };
};
export const bookmarkStore = Bookmarks();

// store to handle the app state
const appState = () => {
  const { subscribe, update } = writable(false);
  return {
    subscribe,
    error: () => update(error => !error),
  };
};

export const AppStore = appState();