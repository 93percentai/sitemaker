#!/usr/bin/env node

import path from "node:path";
import { fileURLToPath } from "node:url";
import { loadConfig } from "./config.js";
import { build } from "./build.js";
import { serve } from "./server.js";
import { scaffold } from "./scaffold.js";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const defaultsDir = path.resolve(__dirname, "..", "defaults");

function main() {
  const args = process.argv.slice(2);
  const command = args[0];

  switch (command) {
    case "build":
      cmdBuild(args.slice(1));
      break;
    case "serve":
      cmdServe(args.slice(1));
      break;
    case "init":
      cmdInit(args.slice(1));
      break;
    case "help":
    case "--help":
    case "-h":
      printUsage();
      break;
    default:
      if (command) {
        console.error(`Unknown command: ${command}\n`);
      }
      printUsage();
      process.exit(1);
  }
}

function cmdBuild(args: string[]) {
  const configPath = getFlag(args, "--config") || "sitemaker.toml";
  const cfg = loadConfig(configPath);

  try {
    build(cfg, defaultsDir);
    console.log("Build complete.");
  } catch (err: any) {
    console.error(`Build failed: ${err.message}`);
    process.exit(1);
  }
}

function cmdServe(args: string[]) {
  const configPath = getFlag(args, "--config") || "sitemaker.toml";
  const port = parseInt(getFlag(args, "--port") || "8080", 10);
  const cfg = loadConfig(configPath);

  serve(cfg, defaultsDir, port).catch((err: any) => {
    console.error(`Server error: ${err.message}`);
    process.exit(1);
  });
}

function cmdInit(args: string[]) {
  const tmpl = getFlag(args, "--template") || "personal";
  const dir = args.find((a) => !a.startsWith("--")) || ".";

  try {
    scaffold(dir, tmpl);
    console.log(`Initialized sitemaker project in ${dir} (template: ${tmpl})`);
  } catch (err: any) {
    console.error(`Init failed: ${err.message}`);
    process.exit(1);
  }
}

function getFlag(args: string[], flag: string): string | undefined {
  const idx = args.indexOf(flag);
  if (idx >= 0 && idx + 1 < args.length) {
    return args[idx + 1];
  }
  return undefined;
}

function printUsage() {
  console.log(`sitemaker - static site generator for personal and org websites

Usage:
  sitemaker <command> [options]

Commands:
  build     Build the site to the output directory
  serve     Start a dev server with live reload
  init      Initialize a new sitemaker project
  help      Show this help message

Options for build/serve:
  --config  Path to config file (default: sitemaker.toml)

Options for serve:
  --port    Port to serve on (default: 8080)

Options for init:
  --template  Template type: personal or org (default: personal)`);
}

main();
