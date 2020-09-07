import Card from './Card.svelte';

export default {
  title: 'Bookie/Card',
  component: Card,
  argTypes: {
    label: { control: 'text' },
  },
};

const Template = ({ onClick, ...args }) => ({
  Component: Card,
  props: args,

});


export const Primary = Template.bind({});
Primary.args = {
  primary: true,
  label: 'LOOOOL',
};
