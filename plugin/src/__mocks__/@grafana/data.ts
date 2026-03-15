export class AppPlugin {
  addPage(config: { title: string; body: unknown; id: string }): AppPlugin {
    return this;
  }
}
