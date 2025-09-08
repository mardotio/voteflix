import {
  type ComponentPropsWithRef,
  type ReactNode,
  forwardRef,
  memo,
} from "react";

import styles from "./icon.module.scss";

export interface IconProps
  extends Omit<ComponentPropsWithRef<"svg">, "stroke"> {
  /** The CSS `fill` color for the icon. */
  fill?: string;
  /** The icon size in pixels. */
  size?: 32 | 24 | 20 | 16 | 12;
}

export const generateIcon = (
  svgContent: ReactNode,
  displayName: string,
  isFilled = false,
) => {
  const Component = (
    { fill = "white", size = 24, className, ...svgProps }: IconProps,
    ref: IconProps["ref"],
  ) => (
    <svg
      className={`${styles[`icon-${size}`]} ${className}`}
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 32 32"
      fill={isFilled ? fill : undefined}
      stroke={isFilled ? undefined : fill}
      ref={ref}
      {...svgProps}
    >
      {svgContent}
    </svg>
  );

  if (import.meta.env.NODE_ENVIRONMENT !== "production") {
    Component.displayName = displayName;
  }

  return memo(forwardRef(Component));
};
