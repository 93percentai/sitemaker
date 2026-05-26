[TITLE]: # (Why I Switched to Neovim in 2026)
[DATE]: # (2026-03-05)
[TAGS]: # (tools, opinion, neovim)
[INCLUDES]: # (H, F)
[INHERITS]: # (post.html)

# Why I Switched to Neovim in 2026

I was a VS Code user for five years. I had the perfect setup — 47 extensions, a custom theme, keybindings tuned to muscle memory. Then my laptop died, and I realized I couldn't reproduce my setup.

That was the push I needed.

## The Breaking Point

Setting up VS Code on a new machine took me **three hours**. Three hours of:

1. Installing extensions one by one
2. Syncing settings that didn't quite sync
3. Fixing conflicts between extensions
4. Remembering which settings I'd changed manually

Meanwhile, my Neovim-using colleague cloned a dotfiles repo and was productive in 5 minutes.

## What Convinced Me

It wasn't the speed, though Neovim is absurdly fast. It wasn't the modal editing, though that grew on me. It was **reproducibility**.

My entire Neovim config is ~200 lines of Lua in a single file[^1]. It lives in my dotfiles repo. I can set up a new machine in under a minute:

```bash
git clone https://github.com/danakim/dotfiles ~/.config
nvim --headless "+Lazy! sync" +qa
```

That's it. Done. Every plugin installed, every keybinding configured, LSP servers downloading in the background.

## My Setup

For the curious, here's my plugin list:

```lua
require("lazy").setup({
    "neovim/nvim-lspconfig",
    "hrsh7th/nvim-cmp",
    "nvim-telescope/telescope.nvim",
    "nvim-treesitter/nvim-treesitter",
    "lewis6991/gitsigns.nvim",
    "folke/tokyonight.nvim",
    "echasnovski/mini.nvim",
})
```

Seven plugins. That's it. Compare that to 47 VS Code extensions.

## The Learning Curve

I won't sugarcoat it — the first two weeks were painful. I was massively slower at everything. But by week three, I was back to normal speed. By month two, I was faster than I'd ever been in VS Code.

The key was **not trying to learn everything at once**. I started with just:

- `i` to enter insert mode
- `Esc` to leave it
- `:w` to save
- `:q` to quit
- `/` to search

Everything else came gradually.

## Do I Miss VS Code?

Honestly? Sometimes. The debugger integration was better. The Git UI was nicer. And some extensions (like the REST client) have no good Neovim equivalent.

But I don't miss the startup time. I don't miss the Electron memory usage. And I definitely don't miss the three-hour setup ritual.

---

[^1]: Okay, it's 247 lines. Close enough.

*This is an opinion piece — your mileage may vary. See [[about|more about me]] and my other posts on [[kubernetes-debugging-checklist|debugging Kubernetes]].*
