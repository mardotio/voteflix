import { guildCreate } from "./events/guildCreate";
import { interactionCreate } from "./events/interactionCreate";
import { ready } from "./events/ready";
import { discordClient } from "./utils/discord";
import { BOT_ENVIRONMENT } from "./utils/environment";

discordClient.once("ready", ready);
discordClient.on("interactionCreate", interactionCreate);
discordClient.on("guildCreate", guildCreate);

(async () =>
  await discordClient.login(BOT_ENVIRONMENT.BOT_DISCORD_BOT_TOKEN))();
