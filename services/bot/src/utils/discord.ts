import Discord from "discord.js";

import { BOT_ENVIRONMENT } from "./environment";

export const discordClient = new Discord.Client({
  intents: [
    Discord.GatewayIntentBits.Guilds,
    Discord.GatewayIntentBits.GuildMessages,
  ],
});

export const discordRest = new Discord.REST({ version: "10" }).setToken(
  BOT_ENVIRONMENT.BOT_DISCORD_BOT_TOKEN,
);
