
import { writable } from 'svelte/store';

const Bookmarks = () => {
  // creates a new writable store populated with some initial data
  const { subscribe, update } = writable([
    { id: '1', title: 'Something', url: 'http://google.com', image: "image", description: "desc", archived: false , deleted: false },
    { id: '2', title: 'Something more', url: 'http://ms.com', image: "image", description: "desc", archived: false , deleted: false },
    { id: '3', title: 'Something else', url: 'http://dn.com', image: "image", description: "desc", archived: false , deleted: false },
    { id: '4', title: 'Something again', url: 'http://lol.com', image: "image", description: "desc", archived: false , deleted: false },
  ]);

  return {
    subscribe,
    // method to archive a task, think of a action with redux or Vuex
    archiveTask: id =>
      update(bookmarks =>
        bookmarks.map(bookmark => (bookmark.id === id ? { ...bookmark, archived: true } : bookmark))
      ),
    // method to archive a bookmark, think of a action with redux or Vuex
    pinTask: id =>
      update(bookmarks =>
        bookmarks.map(bookmark => (bookmark.id === id ? { ...bookmark, deleted: false } : bookmark))
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