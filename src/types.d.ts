declare module "markdown-it-footnote" {
  import MarkdownIt from "markdown-it";
  const plugin: MarkdownIt.PluginSimple;
  export default plugin;
}

declare module "markdown-it-task-lists" {
  import MarkdownIt from "markdown-it";
  const plugin: MarkdownIt.PluginWithOptions<{ enabled?: boolean; label?: boolean }>;
  export default plugin;
}

declare module "toml" {
  export function parse(input: string): any;
}
