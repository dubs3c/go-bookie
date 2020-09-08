import Tag from './Tag.svelte';

export default {
  title: 'Bookie/Tag',
  component: Tag,
  argTypes: {
    value: { control: 'text' },
    deletable: { control: "boolean" },
  },
};

const Template = ({ onClick, ...args }) => ({
  Component: Tag,
  props: args,

});

export const Primary = Template.bind({});
Primary.args = {
  value: 'Tag',
  deletable: true,
};