import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  SlashCommandBuilder,
} from "discord.js";

import { BOT_ENVIRONMENT } from "../utils/environment";
import getCommandName from "../utils/getCommandName";
import { generateLoginJwt } from "../utils/jwt";
import { BotCommand } from "../utils/types";

const getLoginButton = (jwt: string) =>
  new ActionRowBuilder<ButtonBuilder>().addComponents(
    new ButtonBuilder()
      .setStyle(ButtonStyle.Link)
      .setLabel("Login")
      .setURL(encodeURI(`${BOT_ENVIRONMENT.BOT_UI_URL}/login?token=${jwt}`)),
  );

export const login: BotCommand = {
  command: new SlashCommandBuilder()
    .setName(getCommandName("login"))
    .setDescription("Sends you a login link for the voteflix UI"),
  handler: async (interaction) => {
    const { guild } = interaction;

    if (!guild) {
      return;
    }

    const targetMember = guild.members.cache.find(
      (m) => m.user.id === interaction.user.id,
    );

    if (!targetMember) {
      return;
    }

    const jwt = generateLoginJwt({
      sub: targetMember.user.id,
      server: targetMember.guild.id,
      username: targetMember.user.username,
      nickname: targetMember.nickname,
      avatar: targetMember.user.avatar,
    });

    await interaction.reply({
      content: jwt,
      components: [getLoginButton(jwt)],
      flags: ["Ephemeral"],
    });
  },
};
