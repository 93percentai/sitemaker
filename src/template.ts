import fs from "node:fs";
import path from "node:path";
import nunjucks from "nunjucks";
import { Config } from "./config.js";
import { Page, Site } from "./content.js";

export class TemplateEngine {
  private env: nunjucks.Environment;
  private cfg: Config;
  private defaultsDir: string;

  constructor(cfg: Config, defaultsDir: string) {
    this.cfg = cfg;
    this.defaultsDir = defaultsDir;

    const loaders: nunjucks.ILoader[] = [];

    if (fs.existsSync(cfg.build.template_dir)) {
      loaders.push(
        new nunjucks.FileSystemLoader(cfg.build.template_dir, {
          noCache: true,
        })
      );
    }

    const tmplType = cfg.template || "personal";
    const typedDir = path.join(defaultsDir, "templates", tmplType);
    const templatesRoot = path.join(defaultsDir, "templates");

    if (fs.existsSync(typedDir)) {
      loaders.push(new nunjucks.FileSystemLoader(typedDir, { noCache: true }));
    }
    if (fs.existsSync(templatesRoot)) {
      loaders.push(
        new nunjucks.FileSystemLoader(templatesRoot, { noCache: true })
      );
    }

    this.env = new nunjucks.Environment(loaders, { autoescape: false });

    this.env.addFilter("formatDate", (d: Date | null, fmt?: string) => {
      if (!d) return "";
      return d.toISOString().slice(0, 10);
    });

    this.env.addFilter("year", () => new Date().getFullYear());
  }

  render(page: Page, site: Site): string {
    if (page.directives.override) {
      return this.renderOverride(page, site);
    }

    let tmplName = page.directives.inherits;
    if (!tmplName) {
      tmplName = page.isBlogPost ? "post.html" : "base.html";
    }

    const data = this.buildData(page, site);
    return this.env.render(tmplName, data);
  }

  renderList(posts: Page[], site: Site, title: string): string {
    const postData = posts.map((p) => this.pageData(p));
    const data = {
      site: this.siteData(site),
      page: {
        title,
        hasHeader: true,
        hasFooter: true,
        hasTOC: false,
        hasSidebar: false,
      },
      posts: postData,
    };
    return this.env.render("list.html", data);
  }

  private renderOverride(page: Page, site: Site): string {
    const data = this.buildData(page, site);
    return this.env.render(page.directives.override, data);
  }

  private buildData(page: Page, site: Site) {
    return {
      site: this.siteData(site),
      page: this.pageData(page),
    };
  }

  private siteData(site: Site) {
    return {
      title: site.config.site.title,
      description: site.config.site.description,
      baseURL: site.config.site.base_url,
      author: site.config.site.author,
      language: site.config.site.language,
      nav: site.config.site.nav,
      year: new Date().getFullYear(),
    };
  }

  private pageData(page: Page) {
    return {
      title: page.title,
      htmlContent: page.htmlContent,
      tocContent: page.tocContent,
      url: page.url,
      date: page.date,
      dateStr: page.dateStr,
      tags: page.tags,
      readingTime: page.readingTime,
      hasHeader: page.hasHeader,
      hasFooter: page.hasFooter,
      hasTOC: page.hasTOC,
      hasSidebar: page.hasSidebar,
      cssAssets: page.cssAssets,
      jsAssets: page.jsAssets,
      isBlogPost: page.isBlogPost,
      prev: page.prev,
      next: page.next,
    };
  }
}
