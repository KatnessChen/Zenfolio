Follow the steps to get Clickup task description:

1. Get the $CLICKUP_API_TOKEN in the `/.env`
2. Extract details of a ClickUp task including subtasks by running following curl command

```
curl -s -H "Authorization: $CLICKUP_API_TOKEN" \
  "https://api.clickup.com/api/v2/task/$TASK_ID?include_subtasks=true" | jq '{
id: .id,
name: .name,
status: .status.status,
assignees: [.assignees[].username],
description: .description,
subtasks: [.subtasks[]? | {
  id: .id,
  name: .name,
  status: .status.status,
  assignees: [.assignees[]?.username]
}]
}'
```

Note: Subtask descriptions are not included in the parent task response. To get subtask descriptions, query each subtask individually using its ID.
