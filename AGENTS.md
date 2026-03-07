# Your prime directive

You are an agentic AI system committed to the following principles.

- **autonomy:** Whenever you see the need for human intervention or confirmation, also make a proposal how the workflow can be adjusted, so that you can obtain the confirmation on your own, using appropriate tools. If necessary, develope those tools.
- **quality:** Maximize automated testing on all levels (unit testing, integration testing, end-to-end testing). Reduce the need for manual testing and validation wherever you can. Note that this aligns well with the 'autonomy' principle.
- **integration:** Always make sure that the CI pipeline works correctly. If there is a pipeline failure, pause every other task and make fix it.

# Automated UI Testing & Self-Diagnostics

To ensure that the AI agent can autonomously verify frontend behavior and diagnose issues without human intervention, adhere to the following standards:

## 1. Automation Bridge
The application maintains an automation bridge between the Go backend and the React frontend via Wails events.
- **Frontend Listener:** The `App.tsx` component MUST contain a listener for the `"automation:command"` event.
- **Action Extensibility:** When adding new UI features (e.g., a new sidebar, filtering, or modal), you MUST add a corresponding action to this listener (e.g., `inspect_sidebar`, `trigger_filter`) that reports the relevant DOM state or component status back to the Go backend using `LogUIState`.
- **Visibility Verification:** Use `getBoundingClientRect()` and `naturalWidth/Height` in automation commands to verify that elements are not just present in the DOM, but correctly rendered and visible.

## 2. Self-Driving Tests
The application supports a "Self-Test" mode triggered by the environment variable `PHOTOO_SELF_TEST=true`.
- **Verification Logic:** Any significant logic change in the backend that affects the UI should be accompanied by an update to the `runSelfTest()` method in `app.go` or a dedicated diagnostic script.
- **Autonomous Feedback Loop:** Use this mechanism to import test data, trigger UI actions, and capture the resulting `LogUIState` and `LogFrontendError` outputs to verify success.

## 3. Diagnostic Tooling
Maintain and extend the tools in `scripts/diagnose/` and `scripts/check_json/`. These provide a "headless" way to verify backend consistency (database vs. filesystem) and API output formats.

# Development Workflow

Adhere to the workflow described in this file, regardless of programming language.

## Test-Driven Development (TDD)

Follow the canonical TDD workflow as defined by Kent Beck:

### The Five-Step Process

1. **Write a test list** - Before coding, document all expected behavioral variants and edge cases
2. **Write one test** - Convert a single item from the list into an automated test with setup, invocation, and assertions
3. **Make it pass** - Write the minimum code to make the test (and all previous tests) pass
4. **Refactor** - Improve the implementation design while keeping tests green
5. **Repeat** - Continue until the test list is empty

### Red-Green-Refactor Cycle

For each test:
- **Red**: Write a failing test first
- **Green**: Write just enough code to pass
- **Refactor**: Clean up while tests stay green

### Key Principles

- **Never write production code without a failing test first**
- **Interface design happens during test writing** - Think about how behavior is invoked
- **Implementation design happens during refactoring** - Think about how it works internally
- **Test order matters** - Pick tests strategically; the sequence affects both the experience and the result

### Common Mistakes to Avoid

- Writing multiple tests before making any pass
- Skipping the refactor step (leads to messy code)
- Mixing implementation decisions into the test list phase
- Refactoring while making tests pass (separate concerns)
- Copying computed values into expected values
- Over-abstracting before seeing duplication patterns

## Git 

### Local branches

You are one of several human and AI agents working on the same project. Each
agent works in their dedicated git worktree. Each worktree corresponds to one
local branch. You can use the tool `wth` to keep your branch as close to `main`
as possible:

- Commit often
- If technically possible, run all tests and checks that will occur in the CI pipeline, locally, to ensure that the pipeline will not be broken by your commit
- For every commuit, use `wth merge --push .` to merge your branch to `main` and push to the remote repo.
- This will trigger a pipeline run. Make sure the pipeline succeeds. If it does not, fix it.

### Reemote branches

In the remote repository, create branches ONLY in the following cases:

- You have been explicitely instructed to do so by the user
- You are testing the CI pipeline or its components

If a remote branch exists, adhere to the following rules:

- **Minimixe the overall lifetime of the branch**
- **Rebase before merging** - Update with latest trunk changes first
- **Delete immediately after merge**

### Commit messages

Keep commit messages concise.

### Before Every Commit

Make sure that:

1. All tests pass
2. Code quality checks pass
3. No skipped tests
4. Changes are coherent and complete

NEVER commit if tests are skipped. The CI chain checks for skipped tests and fails.
