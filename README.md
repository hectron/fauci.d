# fauci.d

![Fauci showing how getting vaccinated is cool](https://cdn.vidyard.com/thumbnails/7814513/50jFqyrE1bufVRU_JXZct5UCKTudWzcX.gif)

This repository contains a Slack bot that finds COVID-19 vaccines in the United States.

The application registers each vaccine as [a Slack `slash command`](https://api.slack.com/interactivity/slash-commands):

- `/pfizer <postal code>`
- `/moderna <postal code>`
- `/jj <postal code>`

For example, if you wanted to find `Pfizer` vaccinations near `60601`, you'd run the following command:

```
/pfizer 60601
```

## Architecture

This repository uses [Serverless](https://www.serverless.com) to provide bot functionality to Slack. We use
[Serverless](https://serverless.com) to spin up the infrastructure needed to run an AWS lambda (we could also spin up
BCP or Azure  functions using Serverless).

The functions are stored in the `functions` directory, by chat program.

Within the Slack directory, there are two folders -- `backend` and `handler`. The `handler` is the interface that we
provide to the Slack. The `backend` is the interface that gathers the providers, and presents them to the user. We use
this pattern because Slack slash commands need to be acknowledged within 3 seconds. Some provider searches can take over
8 seconds.

The handler receives the request, and immediately responds to Slack with a 200, indicating that
we've acknowledged the request. The handler also asynchronously invokes the `backend` lambda, which can take a little
bit longer to load the providers and notify the user in Slack.

## Development

If you want to test out this application locally, you'll need to [set up a Slack App for your workspace](https://app.slack.com/apps-manage/).  You'll also want to setup the following environment variables:

- `SLACK_API_TOKEN` - the token for your bot
- `MAPBOX_API_TOKEN` - the token to use to communicate with Mapbox
- `SENTRY_DSN` - URL to post messages to sentry
- `SENTRY_ENVIRONMENT` - which environment the application is running on
- `SENTRY_RELEASE` - used to tag the errors in sentry

Useful files:

- `examples/sample_bot.go` -- how to search for a vaccine and post a message to Slack with the results.
- `slack/bot_manifest.yml` -- example configuration settings for a Slack app.

To run the example:

```bash
go run example.go
```

## Running the test suite

- Test entire suite: `go test ./...` or `make test`
- Test individual package: `go test <package>`
- Test vaccines package: `go test github.com/hectron/fauci.d/vaccines`

## Deployment

The application is now automatically deployed using Github Actions when a new **Release** is created. Check the file `.github/workflows/cd.yml` for more information.

If you add any environment variable dependencies, add them to the repository's secrets in Github.

### Manual Deployment

In order to **manually** deploy the application, you will need:

- [Make](https://www.gnu.org/software/make/)
- [serverless CLI](https://www.serverless.com/framework/docs/providers), [**configured to support AWS**](https://www.serverless.com/framework/docs/providers/aws/guide/installation)
- Environment variables required in [`serverless.yml`](https://github.com/hectron/fauci.d/blob/main/serverless.yml#L52-L57) (these can be available using `.envrc`)

### Actual deployment

```bash
make deploy
```
### Testing deploy

```bash
serverless invoke -f slackbot -t
```

### Destroying the lambda (permanently)

```bash
serverless remove
```
