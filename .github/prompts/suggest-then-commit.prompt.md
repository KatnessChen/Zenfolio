# GitHub Copilot Agent Mode Commit Prompt

## Goal:

Instruct GitHub Copilot Agent Mode to generate concise and descriptive Git commit messages. Messages must follow the Conventional Commits specification for a clear project history, easier debugging, and automated changelog generation.

## Desired Commit Pattern:

Commit messages must follow this structure for clarity and automated parsing:

():

[Optional body explaining the change in more detail]

## Guidelines:

- **`<type>`**: Mandatory prefix categorizing the change. Choose only one:
  - `feat`: New feature or enhancement.
  - `fix`: Bug fix or error correction.
  - `docs`: Documentation-only changes.
  - `style`: Formatting or whitespace changes.
  - `refactor`: Code restructuring without functional change.
  - `perf`: Performance improvements.
  - `test`: Adding or correcting tests.
  - `chore`: Routine maintenance, tool updates, dependency upgrades.
  - `build`: Changes affecting the build system or external dependencies.
  - `ci`: CI configuration or script changes.
  - `revert`: Reverts a previous commit.
- **`<scope>` (Optional)**: Brief, descriptive phrase in parentheses indicating the affected codebase area (e.g., `(auth)`, `(UI)`). Omit if global.
- **`<subject>`**: Concise headline (aim for max 50 chars, up to 72 if absolutely necessary), imperative, present tense. Describes _what_ the commit does. Must be in the imperative mood (e.g., 'add', 'fix', not 'added', 'fixes'). Do not capitalize the first letter, and no period at the end.
- **Body (Optional)**: Detailed explanation after a blank line. Describes 'why' and 'how' (problem, motivation, approach, implications). Wrap lines at 72 characters.
- **Avoid issue references**: Do not include (e.g., `Closes #123`); user adds these manually.
- **Analyze staged changes**: Thoroughly examine staged changes to determine the correct `<type>`, `<subject>`, and if a `<body>` is needed.

## Example Generated Commit Message:

feat(auth): add user registration endpoint

This commit introduces a new API endpoint for user registration.
It includes validation for email and password, and hashes the password before saving to the database.

## Execution

1. Execute the `git status` command in the terminal to review staged changes and ensure they are ready for commit.
2. If any changes are not staged, please ignore them and focus only on the staged changes for the commit message.
3. Execute the `git commit -m "your commit message"` in the terminal to create the commit with the suggested message. Ensure the commit message is valid and follows the specified pattern before executing.
