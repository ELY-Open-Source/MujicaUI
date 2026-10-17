# AGENTS.md

MujicaUI: court-gothic components for MyGo native UI. Go library, module `git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI`, one package per component category.
Progress truth: `TASKS.md`. Requirements: `PRD.md` and `docs/prd/ext-*.md` (101–500).
Planned: once 001–500 land, the root package splits into category sub-packages (TASKS T22). Until then everything stays in the root package.

## Stack

| Item | Value |
| --- | --- |
| Go | 1.27.1 (MyGo's minimum) |
| MyGo | `github.com/egoist/mygo v0.5.5` (commit 4d8673c), pinned in go.mod |
| Icons | Lucide 1.52.0, ISC, copied into `icons/svg` |
| Web mirror | `examples/gallery -web :8080`: Go stdlib server drives the Gallery through `ui.NewTester` and streams frames to a canvas. The browser shows pixels; Go runs the components. |

Read MyGo source before using an API: `go env GOMODCACHE`/github.com/egoist/mygo@v0.5.5 (`ui/*.go`, `docs/ui/*.md`, migration guide `docs/ui/migration.md`). Never guess.

## Layout

```text
core/                Use, Tokens, Density, Motion, FontSize, Msg; shared result types and helpers;
                     messages_<group>.go built-in copy (en-US, zh-CN) via addMessages; docs tests
<category>/          one component per file, package = folder name:
                     layout display input datetime navigation data overlay feedback chart
                     finance chat agent git code media project files messaging devtools account
<category>/*_test.go tests beside the component; example_<group>_test.go compiled doc examples
theme/               tokens, density, type, motion, contrast
icons/               embedded Lucide icons: icons.Must("check")
internal/            focus ring, ornaments; internal/snaptest runs snapshot goldens
examples/gallery/    native Gallery; pages_<group>.go per PRD 4.x group
docs/components/     NNN-name.md, one page per component
testdata/snapshots/  <os>/<name>-<light|dark>.png, shared by every package
```

Gallery categories map to packages: Inputs, Choices, Forms, Essentials -> input; Date & Time, Calendar -> datetime; Feedback, Agent Ops & Motion -> feedback; Charts I, Charts II, Market Charts -> chart; Trading, Markets -> finance; Chat Messages, Chat Input, Chat Sessions -> chat; Terminal & Code -> code; Projects -> project; Files & Collaboration -> files; Messaging & Mail -> messaging; Dashboards & Dev Tools -> devtools; Settings & Account -> account. Imports point down only: core imports no category; a piece two categories need moves to core or internal.

## Component rules

- Signature: `func Name(c *ui.Context, value *T, opts NameOptions) ui.Element`. Composites may return a named struct.
- Result structs expose the element as a named field `Element ui.Element`, never embedded; element plus changed is `ChoiceResult`.
- Shared helpers: `iconAction` (borderless icon button), `overlayStyle(c, panel, overlayEnter(c, panel))` (popup surface and entry fade), `internal.RuleDiamond` / `internal.Frame` (ornaments), `colorFade` (state color transition).
- Zero `Options` gives the default look. No `map[string]any` config.
- Explicit font sizes: `FontSize(c, theme.CaptionSize)`; it applies the desktop text scale.
- Colors: `k := core.Tokens(c)`. Metrics: `theme.*` consts and `Density(c).Unit()` / `.ControlHeight()`.
- Build on MyGo Bases (`ui.ButtonBase`, `ui.CheckboxBase`, `ui.SelectBase`, ...) or MyGo widgets restyled through the theme `Use` sets. Never reach into MyGo internals.
- Focus: `internal.Ring(e, k.Focus, radius)` for the 2 DIP ring with 2 DIP gap.
- State marks (checked, selected, progress) use `OnAccent` on `Accent`, or `AccentText`, never `Accent` alone on dark surfaces (2.44:1).
- Ornament color only for dividers, corners, panel titles. No gradients, glow, glass or texture.
- Animated paths: pass coordinates through `internal.Snap(p, v)`. MyGo caches path masks by quarter-pixel keys (ui/path.go); unsnapped motion reuses near-miss masks and breaks snapshots.
- Durations: `core.Motion(c, theme.StateDuration)`; popups use `theme.PopupDuration`.
- Disabled blocks input before handling it (`Disabled(true)`); MyGo dims disabled elements, do not add opacity.
- Copy shown by a component comes from `core.Msg(c, key)`; business copy comes from the caller.
- `ui.Element` is a handle valid only during its build pass: no cross-frame caching, no goroutines touching elements, no element methods inside `HandleInput` or other callbacks (read state such as `IsDisabled` into a local during the build). Keep a `ui.Handle` for persistent identity and `c.Services()` for clipboard or URLs from callbacks. The absent element is `ui.Element{}`; test it with `.Valid()`.
- Overlay content takes initial focus with `core.AutoFocus(c, e)`, not `e.AutoFocus()`: MyGo v0.5.5 records an overlay's opener after its content builds, so same-frame autofocus loses focus restore on close.
- Tests run on the Tester's virtual clock: `tt.Advance(d)` and `tt.Now()`, never `time.Sleep` or `time.Now()`; read element state inside the view.
- Fail fast: panic on programmer errors (bad options); never swallow errors.
- Log state transitions in composites with `slog.Debug("mujica: <component> <event>", ...)`. Never log passwords.
- Files under 500 lines. Comments: one line, only when needed.

## Tests

- `ui.NewTester` drives click, keys, typing. Assert values, not existence.
- Every component: one `snapshot(t, name, w, h, view)` call (light and dark; each package wraps `internal/snaptest`). Refresh with `go test -run <Name> -update ./<category>`
- Run `go vet ./... && go test ./...` before reporting done.
- Each snapshot test renders in a fresh child process (snapshot_test.go), because MyGo's path-mask cache is process-wide. Order and neighbours never change pixels; `-run <Name> -update` is safe.

## Gallery page

Each component adds one `page` to its group slice in `examples/gallery/pages_<group>.go`: ID, Name, ZhName, Category, Demo, Code, Keys. Demo state lives in `ui.Local(c.Root(), "<id>", ...)`. Demos call `log` for every event.

## Review

Each TASKS.md item is reviewed by `codex exec -m gpt-6.1-sol -c model_reasoning_effort=high` (max 3 rounds). This session runs outside Herdr (`HERDR_ENV` unset), so reviews run as direct `codex exec`.
