module.exports = [
  // This is an example redirect, it can be removed once other redirects have been added
  {
    source: "/__test",
    destination: "/",
    permanent: true,
  },
  {
    source: "/__test-local",
    destination: "/",
    permanent: false,
  },
];
