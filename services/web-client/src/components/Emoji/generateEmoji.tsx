import { type ComponentPropsWithRef, type ReactNode, memo } from "react";

import styles from "./Emoji.module.scss";

type EmojiData = Required<ReactNode> | null | undefined;

export interface EmojiProps
  extends Omit<ComponentPropsWithRef<"svg">, "stroke" | "fill"> {
  /** The CSS `fill` color for the icon. */
  fill?: string;
  /** The icon size in pixels. */
  size?: 12 | 16 | 20 | 24 | 32 | 36 | 48;
}

interface GenerateEmojiOptions {
  content: EmojiData;
  displayName: string;
}

export const generateEmoji = ({
  displayName,
  content,
}: GenerateEmojiOptions) => {
  const Component = ({ size = 32, className, ...svgProps }: EmojiProps) => {
    return (
      <span
        className={`${styles.emoji} ${styles[`emoji-${size}`]} ${className ?? ""}`}
        role="img"
        aria-description={displayName}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 36 36"
          preserveAspectRatio="xMidYMid meet"
          {...svgProps}
        >
          {content}
        </svg>
      </span>
    );
  };

  if (import.meta.env.NODE_ENVIRONMENT !== "production") {
    Component.displayName = displayName;
  }

  return memo(Component);
};
