import { type ComponentPropsWithRef, type ReactNode, memo } from "react";

import styles from "./Icon.module.scss";

type IconData = Required<ReactNode> | null | undefined;

type OptionalIfHasContent<
  V extends object,
  F extends IconData,
  O extends IconData,
> = F extends null | undefined
  ? O extends null | undefined
    ? V
    : Partial<V>
  : Partial<V>;

export type IconProps<
  F extends IconData,
  O extends IconData,
> = OptionalIfHasContent<
  {
    iconStyle: F extends null | undefined
      ? O extends null | undefined
        ? never
        : "outline"
      : O extends null | undefined
        ? "solid"
        : "solid" | "outline";
  },
  F,
  O
> & {
  /** The CSS `fill` color for the icon. */
  fill?: string;
  /** The icon size in pixels. */
  size?: 12 | 16 | 20 | 24 | 32 | 36 | 48;
} & Omit<ComponentPropsWithRef<"svg">, "stroke">;

interface GenerateIconOptions<F extends IconData, O extends IconData> {
  filledSvg?: F;
  outlineSvg?: O;
  displayName: string;
}

export const generateIcon = <
  F extends IconData = null,
  O extends IconData = null,
>({
  displayName,
  filledSvg,
  outlineSvg,
}: GenerateIconOptions<F, O>) => {
  const Component = ({
    iconStyle,
    size = 32,
    className,
    ...svgProps
  }: IconProps<F, O>) => {
    if (!filledSvg && !outlineSvg) {
      throw new Error("SVG icon must define at least one style");
    }

    const iS = iconStyle ? iconStyle : outlineSvg ? "outline" : "solid";

    return (
      <span
        className={`${styles.icon} ${styles[`icon-${size}`]} ${className ?? ""}`}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          {...svgProps}
        >
          {iS === "solid" ? filledSvg : outlineSvg}
        </svg>
      </span>
    );
  };

  if (import.meta.env.NODE_ENVIRONMENT !== "production") {
    Component.displayName = displayName;
  }

  return memo(Component);
};
