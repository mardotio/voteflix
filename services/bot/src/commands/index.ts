import { BotCommand } from "../utils/types";
import { login } from "./login";
import { pickMovie } from "./pickMovie";

export const COMMANDS: BotCommand[] = [login, pickMovie];

export const COMMAND_MAP = COMMANDS.reduce<
  Record<string, BotCommand["handler"]>
>(
  (mapped, c) => ({
    ...mapped,
    [c.command.name]: c.handler,
  }),
  {},
);
