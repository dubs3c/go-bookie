import Button from './Button.svelte';

export default {
  title: 'Bookie/Button',
  component: Button,
  argTypes: {
    value: { control: 'text' },
    large: { control: "boolean"},
    fullWidth: { control: "boolean"},
    danger: { control: "boolean"}
  },
};

const Template = ({ onClick, ...args }) => ({
  Component: Button,
  props: args,

});

export const Primary = Template.bind({});
Primary.args = {
  value: 'Submit',
};

export const Danger = Template.bind({});
Danger.args = {
  value: 'Danger',
  danger: true,
};


export const Large = Template.bind({});
Large.args = {
  value: 'Large',
  large: true,
};

export const FullWidth = Template.bind({});
FullWidth.args = {
  value: 'Full Width button',
  large: true,
  fullWidth: true,
};

