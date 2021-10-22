# fauci.d

This repository contains a Slack bot that finds COVID-19 vaccines in the United States.

The application registers each vaccine as [a Slack `slash command`](https://api.slack.com/interactivity/slash-commands):

- `/pfizer <postal code>`
- `/moderna <postal code>`
- `/jj <postal code>`

For example, if you wanted to find `Pfizer` vaccinations near `60601`, you'd run the following command:

```
/pfizer 60601
```

## Development

If you want to test out this application locally, you'll need to [set up a Slack App for your workspace](https://app.slack.com/apps-manage/).  You'll also want to setup the following environment variables:

- `SLACK_API_TOKEN` - the token for your bot
- `MAPBOX_API_TOKEN` - the token to use to communicate with Mapbox

Useful files:

- `example.go` -- how to search for a vaccine and post a message to Slack with the results.
- `slack_bot_manifest.yml` -- example configuration settings for a Slack app.

To run the example:

```bash
go run example.go
```

## Running the test suite

- Test entire suite: `go test ./...`
- Test individual package: `go test <package>`
- Test vaccines package: `go test github.com/hectron/fauci.d/vaccines`

