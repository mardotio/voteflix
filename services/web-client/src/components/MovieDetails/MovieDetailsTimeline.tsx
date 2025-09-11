import { format } from "date-fns";

import type {
  MovieDetailsUsersMap,
  MovieDetailsVote,
  MovieStatus,
} from "../../utils/client/moviesApi";
import { ClapperBoardClosedIcon, LikeIcon } from "../Icon";
import styles from "./MovieDetails.module.scss";

interface MovieDetailsTimelineProps {
  createdAt: number;
  status: MovieStatus;
  creatorId: string;
  votes: MovieDetailsVote[];
  users: MovieDetailsUsersMap;
  watchedAt: number | null;
}

const toTimestamp = (t: number) => format(new Date(t), "EEE, LLL do yyyy, p");

const getVoteLine = (status: MovieStatus, votes: MovieDetailsVote[]) => {
  if (status === "pending") {
    return null;
  }

  const lastVote = votes[votes.length - 1];

  return (
    <li>
      <span className={styles.green} />
      <div className={styles.details}>
        <p className={styles.timestamp}>
          {toTimestamp(lastVote.updatedAt ?? lastVote.createdAt)}
        </p>
        <p className={styles.summary}>
          Reached {status === "rejected" ? "rejection" : "approval"} threshold
        </p>
      </div>
    </li>
  );
};

const getWatchedLine = (status: MovieStatus, watchedAt: number | null) => {
  if (status !== "watched" || !watchedAt) {
    return null;
  }

  return (
    <li>
      <ClapperBoardClosedIcon size={12} />
      <div className={styles.details}>
        <p className={styles.timestamp}>{toTimestamp(watchedAt)}</p>
        <p className={styles.summary}>Marked as watched</p>
      </div>
    </li>
  );
};

export const MovieDetailsTimeline = ({
  status,
  createdAt,
  creatorId,
  votes,
  users,
  watchedAt,
}: MovieDetailsTimelineProps) => {
  return (
    <ul className={styles.timeline}>
      <li>
        <span />
        <div className={styles.details}>
          <p className={styles.timestamp}>{toTimestamp(createdAt)}</p>
          <p className={styles.summary}>Suggested by {users[creatorId].name}</p>
        </div>
      </li>
      {votes.map((v) => (
        <li key={`vote-${v.userId}`}>
          <LikeIcon
            size={12}
            className={v.isApproval ? undefined : styles.invert}
          />
          <div className={styles.details}>
            <p className={styles.timestamp}>{toTimestamp(v.createdAt)}</p>
            <p className={styles.summary}>
              {v.isApproval ? "Approved" : "Rejected"} by {users[v.userId].name}
            </p>
          </div>
        </li>
      ))}
      {getVoteLine(status, votes)}
      {getWatchedLine(status, watchedAt)}
    </ul>
  );
};
