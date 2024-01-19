# ZjeBot

Discord & Twitch chatbot

## Quick Start

To run the bot 2 files need to exist: `secrets.toml` and `data.toml`.

```toml
# secrets.toml
[twitch]
auth = "oauth:oauth_token"
username = "bot_name"
id = "bot_id"

[discord]
auth = "discord_token"
username = "bot_name"
id = "bot_id"
```
```toml
# data.toml
[today]
  Text = "today"

[project]
  Text = "project"
```

Then run the bot:
```shell
$ go run .
```

## License

This project is under the [MIT](./LICENSE) License.
