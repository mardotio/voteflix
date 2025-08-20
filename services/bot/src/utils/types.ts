import {
  ClientEvents,
  CommandInteraction,
  SlashCommandBuilder,
} from "discord.js";

export type EventHandler<T extends keyof ClientEvents> = (
  ...args: ClientEvents[T]
) => void;

export interface BotCommand {
  command: SlashCommandBuilder;
  handler: (interaction: CommandInteraction) => Promise<void> | void;
}
