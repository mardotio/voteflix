import { formatDistanceToNow } from "date-fns";
import { EmbedBuilder, SlashCommandBuilder } from "discord.js";

import { BotApi, PickMovieResponse, isErrorResponse } from "../utils/client";
import getCommandName from "../utils/getCommandName";
import { BotCommand } from "../utils/types";

const createMovieEmbed = (movie: PickMovieResponse) => {
  const formattedDate = formatDistanceToNow(new Date(movie.createdAt), {
    addSuffix: true,
  });

  const embed = new EmbedBuilder()
    .setColor("#0c8be8")
    .setTitle(movie.name)
    .addFields({ name: "Added", value: formattedDate });

  if (movie.creator.avatarUrl) {
    embed.setAuthor({
      name: movie.creator.name,
      iconURL: movie.creator.avatarUrl,
    });
  } else {
    embed.setAuthor({ name: movie.creator.name });
  }
  return embed;
};

export const pickMovie: BotCommand = {
  command: new SlashCommandBuilder()
    .setName(getCommandName("pick-movie"))
    .setDescription("Picks a random pending movie"),
  handler: async (interaction) => {
    const { guild } = interaction;

    if (!guild) {
      return;
    }

    const res = await BotApi.pickMovie(guild.id);

    if (!isErrorResponse(res)) {
      await interaction.reply({ embeds: [createMovieEmbed(res)] });
      return;
    }

    if (res.status === 422) {
      const embed = new EmbedBuilder()
        .setTitle("Pending movies in your list!")
        .setDescription("Get to work and approve some movies")
        .setColor("Red")
        .setImage("https://c.tenor.com/Vyg73kR334sAAAAd/tenor.gif");

      await interaction.reply({ embeds: [embed], flags: ["Ephemeral"] });
      return;
    }

    await interaction.reply("Could not pick a movie at this time");
  },
};
