import Discord, {Guild, GuildMember, PartialUser, User} from "discord.js";
import {BOT_ENVIRONMENT} from "../utils/environment";
import {BotApi, isSuccessResponse} from "../utils/api";

export const createList = async (user: User | PartialUser, guildUser: GuildMember, guild: Guild): Promise<string> => {
    const res = await BotApi.createList({
        discordAvatarId: user.avatar,
        discordNickname: guildUser.nickname,
        discordServerId: guild.id,
        discordServerName: guild.name,
        discordUserId: user.id,
        discordUsername: user.username || ""
    })

    if (isSuccessResponse(res)) {
        return "Welcome! Created a new movie list just for you and your friends"
    }

    if (res.status === 409) {
        return "Glad to have you back! Your movie list is just like you left it"
    }

    return "Failed to create a movie list"
}

export const guildCreate = async (guild: Guild) => {
    try {
        const botAuditLogs = await guild.fetchAuditLogs({
            type: Discord.AuditLogEvent.BotAdd,
            limit: 1,
        });
        const auditEntry = botAuditLogs.entries.first();

        if (
            auditEntry &&
            auditEntry.targetId === BOT_ENVIRONMENT.BOT_DISCORD_APP_ID &&
            auditEntry.executor
        ) {
            const inviter = auditEntry.executor;
            const guildMember = await guild.members.fetch(inviter.id);
            const welcomeMessage = await createList(inviter, guildMember, guild);
            await guild.systemChannel?.send(welcomeMessage);
            return;
        }

        const owner = await guild.fetchOwner();

        const welcomeMessage = await createList(owner.user, owner, guild);

        await guild.systemChannel?.send(welcomeMessage);
    } catch (e) {
        console.error(e);
    }
}