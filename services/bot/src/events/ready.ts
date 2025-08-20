import { Routes } from "discord-api-types/v10";

import { COMMANDS, COMMAND_MAP } from "../commands";
import { discordRest } from "../utils/discord";
import { BOT_ENVIRONMENT } from "../utils/environment";
import { EventHandler } from "../utils/types";

export const ready: EventHandler<"ready"> = async (client) => {
  console.log(`Logged in as ${client.user?.tag}`);
  try {
    console.log("Registering commands");

    await discordRest.put(
      Routes.applicationCommands(BOT_ENVIRONMENT.BOT_DISCORD_APP_ID),
      {
        body: COMMANDS.map((c) => c.command.toJSON()),
      },
    );

    console.log(
      `Finished registering commands:\n  - ${Object.keys(COMMAND_MAP).join("\n  - ")}`,
    );
  } catch (e) {
    console.error("Error registering commands", e);
  }
};
