import Bookmark from './Bookmark.svelte';

export default {
  title: 'Bookie/Bookmark',
  component: Bookmark,
  argTypes: {
    url: { control: 'text' },
    title: { control: "text" },
  },
};

const Template = ({ onClick, ...args }) => ({
  Component: Bookmark,
  props: args,

});

export const Primary = Template.bind({});
Primary.args = {
  url: 'https://dubell.io',
  title: "Hacking the planet!",
};