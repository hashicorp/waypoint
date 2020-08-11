export function setupRouter(container) {
  return container.owner.lookup('router:main').setupRouter();
}
