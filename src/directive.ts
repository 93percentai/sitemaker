export interface Directives {
  includes: string[];
  inherits: string;
  override: string;
  date: string;
  title: string;
  tags: string[];
  draft: boolean;
  slug: string;
}

const DIRECTIVE_RE =
  /^\[(INCLUDES|INHERITS|OVERRIDE|TITLE|DATE|TAGS|DRAFT|SLUG)\]:\s*#\s*\(([^)]*)\)\s*$/i;

export function parseDirectives(content: string): {
  directives: Directives;
  body: string;
} {
  const d: Directives = {
    includes: [],
    inherits: "",
    override: "",
    date: "",
    title: "",
    tags: [],
    draft: false,
    slug: "",
  };

  const lines = content.split("\n");
  let bodyStart = 0;

  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim();
    if (trimmed === "") continue;

    const m = trimmed.match(DIRECTIVE_RE);
    if (!m) {
      bodyStart = i;
      break;
    }

    bodyStart = i + 1;
    const key = m[1].toUpperCase();
    const val = m[2].trim();

    switch (key) {
      case "INCLUDES":
        d.includes.push(...splitCSV(val));
        break;
      case "INHERITS":
        d.inherits = val;
        break;
      case "OVERRIDE":
        d.override = val;
        break;
      case "TITLE":
        d.title = val;
        break;
      case "DATE":
        d.date = val;
        break;
      case "TAGS":
        d.tags.push(...splitCSV(val));
        break;
      case "DRAFT":
        d.draft = val.toLowerCase() === "true";
        break;
      case "SLUG":
        d.slug = val;
        break;
    }
  }

  const body = lines.slice(bodyStart).join("\n");
  return { directives: d, body };
}

export function layoutComponents(d: Directives): string[] {
  return d.includes
    .filter((s) => !s.includes("."))
    .map((s) => s.toUpperCase());
}

export function assetFiles(d: Directives): string[] {
  return d.includes.filter((s) => s.includes("."));
}

export function hasComponent(d: Directives, code: string): boolean {
  return layoutComponents(d).includes(code.toUpperCase());
}

function splitCSV(s: string): string[] {
  return s
    .split(",")
    .map((p) => p.trim())
    .filter((p) => p !== "");
}
