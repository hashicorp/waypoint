import hbs from 'htmlbars-inline-precompile';
import { withKnobs, select, text } from '@storybook/addon-knobs';

export default {
  title: 'Notifications',
  component: 'Notifications',
  decorators: [withKnobs],
};

// add stories by adding more exported functions
export let Notifications = () => ({
  template: hbs`
    <Button {{action this.triggerNotification}}>
      Trigger "{{this.variant}}" Notification
    </Button>
    <Notifications />
  `,
  context: {
    variant: select('@variant', ['success', 'info', 'warning', 'error'], 'info'),
    title: text('Title', 'My Title'),
    content: text('Content', 'My Content'),
    actionText: text('Action Text', 'Action!'),
    triggerNotification: function() {
      this.flashMessages[this.variant](this.title, {
        content: `${this.content}`,
        actionText: `${this.actionText}`,
        onAction: function() {
          alert('hi');
        }
      });
    },
  }
});
