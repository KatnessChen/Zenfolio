# GitHub Copilot Agent Mode Branch Name Prompt

## Goal:

Instruct GitHub Copilot Agent Mode to generate a concise, descriptive Git branch name based on task/issue content, and automatically provide the `git checkout -b` command for a streamlined workflow.

## Desired Branch Name Pattern:

Branch names should follow this structure for easy identification and automation:

**With Task ID**: `<type>/<issue-id>/<short-description>`
**Without Task ID**: `<type>/<short-description>`

## Desired Output:

The final output should be the generated branch name, followed by the `git checkout -b` command using that name, formatted as follows:

Suggested Branch Name: Command to create branch: git checkout -b

## Guidelines:

To ensure consistent and useful branch names, adhere to these guidelines:

- **`<type>`**: Mandatory prefix indicating task type. Choose one:
  - `feat`: New feature.
  - `fix`: Bug fix.
  - `docs`: Documentation changes.
  - `chore`: Maintenance, tool/dependency updates.
  - `refactor`: Code restructuring.
  - `perf`: Performance optimization.
- **`<issue-id>`**: Optional task/issue number (e.g., `JIRA-123`, `GH-456`). Only include if a specific task ID is provided by the user. If no task ID is provided, omit this part entirely.
- **`<short-description>`**: Mandatory brief summary of task content, converted to lowercase, hyphen-separated words (e.g., `implement-user-registration`). Be concise.
- **No Special Characters**: Only lowercase letters, numbers, and hyphens (`-`). No spaces, underscores, or other special characters.
- **Length**: Reasonable length (max 50 chars) for readability.
- **Clarity**: Must clearly represent the task/issue without viewing full content.
- **Command Generation**: The `git checkout -b` command _must_ be generated using the exact branch name produced.

## Example Generated Output:

Here are examples of the desired output based on the above guidelines:

**With Task ID**:
Suggested Branch Name: feat/GH-101/add-user-profile-editing
Command to create branch: git checkout -b feat/GH-101/add-user-profile-editing

**Without Task ID**:
Suggested Branch Name: feat/add-user-profile-editing
Command to create branch: git checkout -b feat/add-user-profile-editing

## Execution:

1. **Task Confirmation**: Always confirm the task details with the user before proceeding. Ask for clarification if the task description is unclear or if a task ID should be included.
2. **Infer from Staged Files**: If the user did not provide a task ID or clear task description, analyze the staged files using `git status` and `git diff --staged` to infer the appropriate branch name based on the changes being made.
3. Execute the `git checkout main` command in the terminal to ensure you're on the main branch before creating a new branch.
4. Execute the `git checkout -b` in the terminal to create the branch with the suggested name. Ensure the branch name is valid and follows the specified pattern before executing.
5. Execute the `git push -u origin <branch-name>` command in the terminal to push the new branch to the origin remote repository and set up upstream tracking.
