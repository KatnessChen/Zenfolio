Extract details of a ClickUp task running following curl command:

```
curl -s -H "Authorization: $CLICKUP_API_TOKEN" \
  "https://api.clickup.com/api/v2/task/$TASK_ID" | jq '{
id: .id,
name: .name,
status: .status.status,
assignees: [.assignees[].username],
description: .description
}'
```

Don't create shell script. Just run the curl command and return the output.
