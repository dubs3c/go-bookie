
import BookmarkList from './BookmarkList.svelte';
import { bookmarkData, actionsData } from './Bookmark.stories';
export default {
  title: 'Bookie/BookmarkList',
  excludeStories: /.*Data$/,
};

export const defaultBookmarksData = [
  { ...bookmarkData, id: '1', title: 'Such cool article', url: 'https://dn.se', image: "image", description: "desc", archived: false , deleted: false },
  { ...bookmarkData, id: '2', title: 'Don\'t be evil, am i right?', url: 'https://google.com', image: "image", description: "desc", archived: false , deleted: false},
  { ...bookmarkData, id: '3', title: 'The Mentor, hacker #1', url: 'https://dubell.io', image: "image", description: "desc", archived: false , deleted: false },
  { ...bookmarkData, id: '4', title: 'WE LOVE LINUX OMG', url: 'https://microsoft.com', image: "image", description: "desc", archived: false , deleted: false },
  { ...bookmarkData, id: '5', title: 'These are not the droids', url: 'https://localhost.com', image: "image", description: "desc", archived: false , deleted: false },
  { ...bookmarkData, id: '6', title: 'Funny cat videos', url: 'https://youtube.com', image: "image", description: "desc", archived: false , deleted: false },
];
export const deletedBookmarksData = [
  ...defaultBookmarksData.slice(0, 5),
  { id: '6', title: 'Funny cat videos', url: 'https://youtube.com', image: "image", description: "desc", archived: false , deleted: true },
];

// default TaskList state
export const Default = () => ({
  Component: BookmarkList,
  props: {
    bookmarks: defaultBookmarksData,
  },
  on: {
    ...actionsData,
  },
});
// tasklist with pinned tasks
export const DeletedBookmarks = () => ({
  Component: BookmarkList,
  props: {
    bookmarks: deletedBookmarksData,
  },
  on: {
    ...actionsData,
  },
});
// BookmarkList in loading state
export const Loading = () => ({
  Component: BookmarkList,
  props: {
    loading: true,
  },
});
// BookmarkList no tasks
export const Empty = () => ({
  Component: BookmarkList,
});