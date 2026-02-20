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
