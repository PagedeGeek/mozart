# mozart
Mozart: simple scheduler

## Run

    $ go run src/github.com/PagedeGeek/mozart/commands/mozart.go

## Use

Create task:

    $ curl -i -X POST http://localhost:1357/tasks/schedule --data @./correct_task.json

correct_task.json
```json
{
  "in": "10s",
  "do": "http_request",
  "timeout": "15s",
  "params": {
    "url": "http://localhost:4000/task_executed",
    "verb": "POST",
    "header": { "X-Auth-Token": "MY_TOKEN" },
    "json_body": { "foo": "bar", "number": 123 }
  }
}
```

Count tasks:

    $ curl -i http://localhost:1357/tasks/count

List tasks:

    $ curl -i http://localhost:1357/tasks

Remove task:

    $ curl -i -X DELETE http://localhost:1357/tasks/unschedule/e4fcbde6-8abd-4a32-865f-885376d80bc6

You can read files:
- mozart_client.rb
- test_mozart_client.Rb
