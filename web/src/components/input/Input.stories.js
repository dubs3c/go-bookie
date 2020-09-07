import Input from './Input.svelte';

export default {
  title: 'Bookie/Input',
  component: Input,
  argTypes: {
    placeholder: { control: 'text' },
    name: { control: 'text' },
    readOnly: { control: 'boolean' },
  },
};

const Template = ({ onClick, ...args }) => ({
  Component: Input,
  props: args,

});

export const Primary = Template.bind({});
Primary.args = {
  placeholder: 'placeholder',
  name: 'form',
  readOnly: false,
};
