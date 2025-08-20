import { BotCommand } from "../utils/types";
import { login } from "./login";

export const COMMANDS: BotCommand[] = [login];

export const COMMAND_MAP = COMMANDS.reduce<
  Record<string, BotCommand["handler"]>
>(
  (mapped, c) => ({
    ...mapped,
    [c.command.name]: c.handler,
  }),
  {},
);
