# YTTG: yt-dlp audio to Telegram
Tool to send audio from online video to Telegram Channel using [yt-dlp](https://github.com/yt-dlp/yt-dlp) with Web UI. See supported sites [here](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md).

## Prerequisites

- Obtain a Telegram [Bot Token](https://core.telegram.org/bots/tutorial#obtain-your-bot-token).
- Add Telegram Bot into Telegram Channel as admin.
- Obtain Telegram [App](https://my.telegram.org/apps) `api_id` and `api_hash`.

## Set up
Copy the `.env.example` file in the root directory of the project to the `.env` file.
Set vars.

| Environment         | Example Value                  | Description               |
|---------------------|--------------------------------|---------------------------|
| TELEGRAM_TOKEN      |                                | telegram token            |
| TELEGRAM_CHAN       |                                | telegram channel          |
| TELEGRAM_TIMEOUT    | `5`                            | telegram timeout, minutes |
| TELEGRAM_SERVER     | `http://telegram-bot-api:8081` | telegram bot api server   |
| TELEGRAM_API_ID     |                                | telegram app api_id       |
| TELEGRAM_API_HASH   |                                | telegram app api_hash     |

## Start

0. Before the very first launch, you need to run
```sh
docker compose build
```
1. Up containers
```sh
docker compose up -d
```
2. Open [http://localhost:8080](http://localhost:8080).
3. Down containers
```sh
docker compose down
```