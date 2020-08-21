import hbs from 'htmlbars-inline-precompile';

const CONFIG = {
  title: 'Components / Error Message',
  component: 'PdsErrorMessage',
};

// add stories by adding more exported functions
const index = () => ({
  template: hbs`
    <Pds::ErrorMessage>
      The quick brown fox jumps over the lazy dog.
    </Pds::ErrorMessage>
  `,
});

const multiple_lines = () => ({
  template: hbs`
    <Pds::ErrorMessage>
      We must neatly release the fail-fast XP complexity! Testing the
      requirements should allow our item to design the steady scrum against the
      metric. Documenting the models should allow our feature to prevent the
      rapid spec toward the iteration. If we release past the tasks, we can get
      the minimum test on the WIP sprint review. If it takes you too long to
      detail the domain MVP product backlog, then you are not refactoring enough.
      Try to test the fibonacci product vision, maybe it will help release the
      WIP development practices. It was discovered that by neatly breaking the
      XP items, we can adapt the pirate VOC trend except the weekly production
      release. If it takes you too long to release the iterative WIP epic, then
      you are not joining enough. Committing the teams should allow our test to
      select the certified sprint planning meeting without the requirement.
      Given aggressive bottlenecks, easily designing the standup product owner
      among the extreme deadlines will usually release the backlog product owner.
    </Pds::ErrorMessage>
  `,
});

export {
  CONFIG as default,

  index,
  multiple_lines,
}
