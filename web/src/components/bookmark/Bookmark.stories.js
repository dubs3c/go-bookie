import Bookmark from './Bookmark.svelte';
import { action } from '@storybook/addon-actions';

export default {
  title: 'Bookie/Bookmark',
  excludeStories: /.*Data$/,
};

export const actionsData = {
  onPinTask: action('onPinTask'),
  onArchiveTask: action('onArchiveTask'),
};

export const bookmarkData = {
  id: '1',
  title: 'Super cool website',
  url: 'https://infected-database.com',
  description: "cool desc",
  image: "image",
  archived: false,
  deleted: false
};

// default task state
export const Default = () => ({
  Component: Bookmark,
  props: {
    bookmark: bookmarkData,
  },
  on: {
    ...actionsData,
  },
});

export const Deleted = () => ({
  Component: Bookmark,
  props: {
    bookmark: {
      ...bookmarkData,
      deleted: true,
      archived: false,
    },
  },
  on: {
    ...actionsData,
  },
});

export const Archived = () => ({
  Component: Bookmark,
  props: {
    bookmark: {
      ...bookmarkData,
      archived: true,
      deleted: false,
    },
  },
  on: {
    ...actionsData,
  },
});