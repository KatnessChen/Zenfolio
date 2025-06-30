# GitHub Copilot Agent Mode Commit Prompt

## Goal:

Instruct GitHub Copilot Agent Mode to generate concise and descriptive Git commit messages. Messages must follow the Conventional Commits specification for a clear project history, easier debugging, and automated changelog generation.

## Desired Commit Pattern:

Commit messages must follow this structure for clarity and automated parsing:

():

[Optional body explaining the change in more detail]

## Guidelines:

- **`<type>`**: One of: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`, `build`, `ci`, `revert`.
- **`<scope>` (Optional)**: Brief phrase for affected area, e.g., `(auth)`, `(UI)`.
- **`<subject>`**: Imperative, present tense, max 50 chars (72 if needed), no period.
- **Body (Optional)**: Explains 'why' and 'how', wrapped at 72 chars.
- **No issue refs**: Do not include (e.g., `Closes #123`).
- **Analyze staged changes**: Use only staged changes to determine type/subject/body.

## Example:

feat(auth): add user registration endpoint

This commit introduces a new API endpoint for user registration.
It includes validation for email and password, and hashes the password before saving to the database.

## Execution

1. Run `git status` to review staged changes.
2. Ignore unstaged changes for the commit message.
3. **Pre-Commit CI Checks**:
   - For `backend/` changes: review and run steps from `.github/workflows/backend-ci.yml`.
   - For `frontend/` changes: review and run steps from `.github/workflows/frontend-ci.yml`.
   - Ensure all checks pass before proceeding.
4. Commit: `git commit -m "your commit message"` (must follow pattern).
5. Push: `git push`.
