# GitHub Pull Request Comments Retrieval Prompt

## Goal:

Instruct GitHub Copilot Agent Mode to retrieve and display all comments from a specific GitHub Pull Request, including review comments, general comments, and suggestions.

## Instructions:

When the user provides a Pull Request ID (e.g., #123), perform the following actions:

1. **Collect All Comments**:

Run commands `gh pr view $PR_ID --comments` to get the comments associated with the Pull Request. Ensure to include:

- **Issue Comments**: General comments on the PR conversation
- **Review Comments**: Line-specific comments from code reviews
- **Review Summaries**: Overall review feedback (approve, request changes, comment)
- **Suggested Changes**: Code suggestions and their status (applied/pending)

2. **Display Format**:
   Present the information in a structured format:

   ````
   # Pull Request #[ID]: [Title]

   **Author**: [Username]
   **Created**: [Date]
   **Status**: [Open/Closed/Merged]
   **Branch**: [source] → [target]

   ## Description
   [PR Description]

   ## Comments ([count] total)

   ### General Comments
   **[Username]** - [Date]
   [Comment content]

   ### Review Comments
   **[Username]** - [Date] - [File]:[Line]
   [Review comment content]

   ### Code Suggestions
   **[Username]** - [Date] - [File]:[Line]
   ```suggestion
   [Suggested code]
   ````

   [Status: Applied/Pending]

   ## Review Summary

   - ✅ **[Username]**: Approved
   - 🔄 **[Username]**: Requested Changes
   - 💬 **[Username]**: Commented

   ```

   ```

3. **Error Handling**:
   - If PR ID doesn't exist, inform the user
   - If repository access is restricted, explain the limitation
   - If no comments exist, display "No comments found"

## Usage Examples:

- "Get comments for PR #42"
- "Show me all comments from pull request #123"
- "Retrieve PR #456 review feedback"
