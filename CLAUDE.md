## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## 5. Bubble Tea v2 Patterns

**Import paths:** `charm.land/bubbletea/v2` and `charm.land/bubbles/v2` (vanity domains, not `github.com/charmbracelet/`).

**Key differences from v1:**
- `tea.KeyMsg` is now an interface; use `tea.KeyPressMsg` for key presses
- `msg.Type` → `msg.Key().Code` (returns a rune constant like `tea.KeyEnter`)
- `tea.KeyCtrlC` doesn't exist — check `key.Code == 'c' && key.Mod == tea.ModCtrl`
- `View()` returns `tea.View`, not `string` — use `tea.NewView("content")`
- `textinput.Model.Update()` must be called with key messages for typing to work
- Use `textinput.Model.Focus()` / `Blur()` to activate/deactivate input (not cursor toggling)
- `textinput.Model.Reset()` clears the input value

**Common pitfall:** If text input doesn't respond to typing, the `Update` method is likely not forwarding `tea.KeyPressMsg` to the focused `textinput.Model.Update()`. Always check input forwarding first.
