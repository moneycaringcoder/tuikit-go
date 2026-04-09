# Changelog

Full release history and changelogs are published with each GitHub Release:

**[github.com/moneycaringcoder/tuikit-go/releases](https://github.com/moneycaringcoder/tuikit-go/releases)**

## Recent Releases

### v0.7.0
- Table detail bars (`DetailFunc`, `DetailRenderer`, `DetailHeight`)
- `RowStyler` full-row background styling with `NoRowStyle` column exemption
- ListView cursor tween animation (120ms ease-out)
- `tuitest` virtual-terminal testing framework with vitest-style reporter
- CLI primitives rendering inline with colored completion states

### v0.6.0
- tuitest framework: 30+ assertions, golden files, session record/replay
- JUnit + HTML reporters
- vitest-like console runner

### v0.5.0
- CLI primitives: Confirm, SelectOne, MultiSelect, Input, Password, Spinner, Progress
- Styled message helpers (Title, Step, Success, Warning, Error, Info)

### v0.4.0
- Tabs component (horizontal + vertical orientation)
- Picker fuzzy command palette with preview pane
- Toast notification system

### v0.3.0
- LogViewer with level filtering, substring search, auto-scroll pause/resume
- ConfigEditor declarative settings overlay
- CommandBar inline command input

### v0.2.0
- Table virtual mode (`TableRowProvider`) for millions of rows
- Table filter callbacks for virtual mode (`OnFilterChange`)
- Theme presets: Dracula, Catppuccin Mocha, Tokyo Night, Nord, Gruvbox Dark, Rose Pine, Kanagawa, One Dark

### v0.1.0
- Initial release: Table, ListView, StatusBar, Help, DualPane layout
- Theme system with semantic color tokens
- Self-update system with SHA256 verification and rollback
- Keybinding registry with auto-generated Help overlay
