import Discord from "discord.js";
import {BOT_ENVIRONMENT} from "./utils/environment";
import {guildCreate} from "./events/guildCreate";

const client = new Discord.Client({
    intents: [
        Discord.GatewayIntentBits.Guilds,
        Discord.GatewayIntentBits.GuildMessages,
    ],
});

client.once("ready", () => {
    console.log(`Logged in as ${client.user?.tag}`);
});

client.on("guildCreate", guildCreate);

(async () => await client.login(BOT_ENVIRONMENT.BOT_DISCORD_BOT_TOKEN))();
