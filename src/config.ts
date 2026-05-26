import fs from "node:fs";
import TOML from "toml";

export interface NavItem {
  label: string;
  url: string;
}

export interface SiteConfig {
  title: string;
  description: string;
  base_url: string;
  author: string;
  language: string;
  nav: NavItem[];
}

export interface BuildConfig {
  content_dir: string;
  template_dir: string;
  static_dir: string;
  output_dir: string;
}

export interface BlogConfig {
  posts_dir: string;
  date_format: string;
  posts_per_page: number;
  enable_rss: boolean;
}

export interface URLConfig {
  style: "flat" | "directory";
  strip_date: boolean;
}

export interface Config {
  site: SiteConfig;
  build: BuildConfig;
  blog: BlogConfig;
  urls: URLConfig;
  template: string;
}

export function defaultConfig(): Config {
  return {
    site: {
      title: "My Site",
      description: "A site built with sitemaker",
      base_url: "/",
      author: "",
      language: "en",
      nav: [],
    },
    build: {
      content_dir: "content",
      template_dir: "templates",
      static_dir: "static",
      output_dir: "dist",
    },
    blog: {
      posts_dir: "posts",
      date_format: "YYYY-MM-DD",
      posts_per_page: 10,
      enable_rss: true,
    },
    urls: {
      style: "flat",
      strip_date: true,
    },
    template: "personal",
  };
}

export function loadConfig(path: string): Config {
  const defaults = defaultConfig();

  if (!fs.existsSync(path)) {
    return defaults;
  }

  const raw = fs.readFileSync(path, "utf-8");
  const parsed = TOML.parse(raw);

  return deepMerge(defaults, parsed) as Config;
}

function deepMerge(target: any, source: any): any {
  const result = { ...target };
  for (const key of Object.keys(source)) {
    if (
      source[key] &&
      typeof source[key] === "object" &&
      !Array.isArray(source[key]) &&
      target[key] &&
      typeof target[key] === "object"
    ) {
      result[key] = deepMerge(target[key], source[key]);
    } else {
      result[key] = source[key];
    }
  }
  return result;
}
