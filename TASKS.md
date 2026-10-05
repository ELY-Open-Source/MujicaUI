# TASKS

Ordered by PRD sections. Check an item only after its codex review passes. Review rounds noted as `r1`, `r2`, `r3`.

- [x] T01 §1–2 Module, MyGo v0.2.10 pin, reuse rules, AGENTS.md — r1 FAIL (line height, enum checks), r2 PASS
- [x] T02 §3 Tokens, density, type, motion, contrast tests, `Use`, icons, i18n — r1 FAIL, r2 PASS
- [x] T03 §4.1 Layout 001–010 — r1 FAIL 8, r2 FAIL 2, r3 FAIL 1 (folded focus); fixed b34958b, verified by coordinator red/green after cap
- [x] T04 §4.2 Display 011–020 — r1 FAIL 5, r2 FAIL 5, r3 FAIL 2 (image sizing, Retry clip); fixed e922651, verified by coordinator tests and snapshot after cap
- [x] T05 §4.3 Inputs 021–030 — r1 FAIL 10, r2 FAIL 5, r3 FAIL 2 (float grid, password in test output); fixed 07a90f4, verified by coordinator red/green after cap
- [ ] T06 §4.4 Choices 031–040
- [ ] T07 §4.5 Forms 041–050 — r1 FAIL 7, r2 FAIL 4, r3 FAIL 1 (disabled cut); fixed 49b3b78, coordinator red/green verified; checks after T06
- [ ] T08 §4.6 Date and time 051–058 — r1 FAIL 9, r2 FAIL 5, r3 FAIL 2 (popup pick commit, selection mark); fixed e3eaef7, coordinator red/green verified; checks after T07
- [ ] T09 §4.7 Navigation 059–068 — r1 FAIL 5, r2 FAIL 5, r3 FAIL 1 (Tab out of menubar); fixed 7e4d862, coordinator red/green verified; checks after T08
- [ ] T10 §4.8 Content and data 069–080 — r1 FAIL 9, r2 FAIL 4, r3 PASS; checks after T09
- [ ] T11 §4.9 Overlays and menus 081–090 — r1 FAIL 12, r2 FAIL 5, r3 PASS; checks after T10
- [ ] T12 §4.10 Feedback 091–100 — r1 FAIL 5, r2 FAIL 4, r3 PASS; checks after T11
- [ ] T13 §4.11 Upstream mapping documented
- [ ] T14 §5 States, keyboard, a11y, input, overlays: cross-cutting tests
- [ ] T15 §6 API and lifecycle audit; README, LICENSE, THIRD_PARTY_NOTICES
- [ ] T16 §7 Platform table, offline, resources
- [ ] T17 §8 Native Gallery: categories, zh/en search, theme, density, code copy, three composite demos
- [ ] T18 §8 Web mirror of the Gallery (user request)
- [ ] T19 §8 Docs: install, minimal example, theme, compat; 100 component pages; examples compile
- [ ] T20 §9 Snapshots, density baselines, large-data and lifecycle tests
- [ ] T21 §9 CI workflow
- [ ] T22 §10 Acceptance checklist
