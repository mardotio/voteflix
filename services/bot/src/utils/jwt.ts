import jwt from "jsonwebtoken";

import { BOT_ENVIRONMENT } from "./environment";

interface JwtPayload {
  sub: string;
  server: string;
  username: string;
  avatar?: string | null;
  nickname?: string | null;
  iat: number;
  exp: number;
}

const SECONDS_IN_MINUTE = 60;
const EXPIRATION_IN_SECONDS = SECONDS_IN_MINUTE * 2;

type JwtCreateRequest = Omit<JwtPayload, "exp" | "iat">;

const getJwtPayload = (user: JwtCreateRequest): JwtPayload => {
  const iat = Math.floor(Date.now() / 1000);
  return {
    ...user,
    iat,
    exp: iat + EXPIRATION_IN_SECONDS,
  };
};

export const generateLoginJwt = (user: JwtCreateRequest) =>
  jwt.sign(getJwtPayload(user), BOT_ENVIRONMENT.BOT_JWT_SECRET);
