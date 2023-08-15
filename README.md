# ChatGPT Telegram bot

This is a telegram bot that allows you to interact with the ChatGPT.

## Usage

First of all, you will need a [Telegram Bot token](https://core.telegram.org/bots#6-botfather)
and an [OpenAI API key](https://beta.openai.com/account/api-keys).

Then you need to copy the `.env.sample` file to `.env` and fill it with your credentials.

You can also specify a list of Telegram user IDs.
In this case, only the specified users will be allowed to interact with the bot.
To get your user ID, message `@userinfobot` on Telegram.
Multiple IDs can be provided, separated by commas.

Run bot:

```sh
go run cmd/server/server.go
```

The bot maintains the conversation context for a specific period of time.
By default, this value is set to 15 minutes from the arrival of the last message, but it can be changed using the
MESSAGES_RETENTION_PERIOD environment variable, which is specified in minutes.
The command /reset can be used to forcibly reset the context.

## Deploy to VDS

To deploy a Golang program to a VDS (Virtual Dedicated Server), you can follow these steps:

1. Build your golang program - Before deploying your program, you need to build it.
   To do that, first, make sure you have golang installed on your local machine.
   Then, navigate to your project directory and run the following command:

```
go build -o go-chatgpt-telegram-bot cmd/server/server.go
```

This will create an executable file in your project directory.

2. Transfer executable file to your VDS. You can do this using an FTP client or SCP. SCP example:

```bash
scp go-chatgpt-telegram-bot deploy@127.0.0.1:~
```

3. Create a new systemd file:

```bash
sudo nano /etc/systemd/system/go-chatgpt-telegram-bot.service
```

Here is an example systemd service file for a Go program:

```
[Unit]
Description=ChatGPT Telegram bot
After=network.target

[Service]
User=deploy
ExecStart=/home/deploy/go-chatgpt-telegram-bot
Environment=TELEGRAM_BOT_TOKEN=secret
Environment=OPENAI_API_KEY=secret
Environment=TELEGRAM_USER_IDS=

[Install]
WantedBy=multi-user.target
```

_Make sure to replace the environment variables with your own values._

4. Save the file and reload systemd:

```bash
sudo systemctl daemon-reload
```

5. Start the service:

```bash
sudo systemctl start go-chatgpt-telegram-bot
```

6. Verify that the service is running:

```bash
sudo systemctl status go-chatgpt-telegram-bot
```

7. If everything is working correctly, enable the service to start at boot:

```bash
sudo systemctl enable go-chatgpt-telegram-bot
```

That's it! Your Golang program should now be running on your VDS and will automatically start up whenever your server
reboots.