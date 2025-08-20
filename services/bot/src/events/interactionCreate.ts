import { COMMAND_MAP } from "../commands";
import { EventHandler } from "../utils/types";

export const interactionCreate: EventHandler<"interactionCreate"> = async (
  interaction,
) => {
  if (!interaction.isCommand()) {
    return;
  }

  const commandHandler = COMMAND_MAP[interaction.commandName];

  if (!commandHandler) {
    return;
  }

  await commandHandler(interaction);
};
