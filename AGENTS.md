# AGENTS.md

MujicaUI: court-gothic components for MyGo native UI. Go library, package `mujica`, module `mujicaui`.
Progress truth: `TASKS.md`. Requirements: `PRD.md`.

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 (MyGo's minimum) |
| MyGo | `github.com/egoist/mygo v0.2.10` (commit 459b512), pinned in go.mod |
| Icons | Lucide 1.52.0, ISC, copied into `icons/svg` |
| Web mirror | `examples/gallery -web :8080`: Go stdlib server drives the Gallery through `ui.NewTester` and streams frames to a canvas. The browser shows pixels; Go runs the components. |

Read MyGo source before using an API: `go env GOMODCACHE`/github.com/egoist/mygo@v0.2.10 (`ui/*.go`, `docs/ui/*.md`). Never guess.

## Layout

```text
*.go                 one component per file, root package mujica
*_test.go            tests beside the component
messages_<group>.go  built-in copy (en-US, zh-CN) via addMessages
theme/               tokens, density, type, motion, contrast
icons/               embedded Lucide icons: icons.Must("check")
internal/            focus ring, ornaments
examples/gallery/    native Gallery; pages_<group>.go per PRD 4.x group
docs/components/     NNN-name.md, one page per component
testdata/snapshots/  <os>/<name>-<light|dark>.png
```

## Component rules

- Signature: `func Name(c *ui.Context, value *T, opts NameOptions) *ui.Element`. Composites may return a named struct.
- Zero `Options` gives the default look. No `map[string]any` config.
- Colors: `k := mujica.Tokens(c)`. Metrics: `theme.*` consts and `Density(c).Unit()` / `.ControlHeight()`.
- Build on MyGo Bases (`ui.ButtonBase`, `ui.CheckboxBase`, `ui.SelectBase`, ...) or MyGo widgets restyled through the theme `Use` sets. Never reach into MyGo internals.
- Focus: `internal.Ring(e, k.Focus, radius)` for the 2 DIP ring with 2 DIP gap.
- State marks (checked, selected, progress) use `OnAccent` on `Accent`, or `AccentText`, never `Accent` alone on dark surfaces (2.44:1).
- Ornament color only for dividers, corners, panel titles. No gradients, glow, glass or texture.
- Animated paths: pass coordinates through `internal.Snap(p, v)`. MyGo caches path masks by quarter-pixel keys (ui/path.go); unsnapped motion reuses near-miss masks and breaks snapshots.
- Durations: `mujica.Motion(c, theme.StateDuration)`; popups use `theme.PopupDuration`.
- Disabled blocks input before handling it (`Disabled(true)`); MyGo dims disabled elements, do not add opacity.
- Copy shown by a component comes from `msg(c, key)`; business copy comes from the caller.
- No cross-frame `*ui.Element` caching, no goroutines touching elements.
- Fail fast: panic on programmer errors (bad options); never swallow errors.
- Log state transitions in composites with `slog.Debug("mujica: <component> <event>", ...)`. Never log passwords.
- Files under 500 lines. Comments: one line, only when needed.

## Tests

- `ui.NewTester` drives click, keys, typing. Assert values, not existence.
- Every component: one `snapshot(t, name, w, h, view)` call (light and dark). Refresh with `go test -run <Name> -update .`
- Run `go vet ./... && go test ./...` before reporting done.

## Gallery page

Each component adds one `page` to its group slice in `examples/gallery/pages_<group>.go`: ID, Name, ZhName, Category, Demo, Code, Keys. Demo state lives in `ui.Local(c.Root(), "<id>", ...)`. Demos call `log` for every event.

## Review

Each TASKS.md item is reviewed by `codex exec -m gpt-6.1-sol -c model_reasoning_effort=high` (max 3 rounds). This session runs outside Herdr (`HERDR_ENV` unset), so reviews run as direct `codex exec`.
