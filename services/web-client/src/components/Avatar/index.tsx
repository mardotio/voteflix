import styles from "./Avatar.module.scss";

interface AvatarProps {
  size?: 12 | 24 | 36 | 48;
  src: string | null;
  name: string;
}

export const Avatar = ({ size = 48, name, src }: AvatarProps) => {
  const [i1, i2] = name
    .replaceAll(/[^a-zA-Z0-9-]/g, "")
    .split(" ")
    .map((v) => v.charAt(0).toLocaleUpperCase());

  if (!src) {
    return (
      <div
        className={`${styles.wrapper} ${styles.empty} ${styles[`size-${size}`]}`}
        role="img"
        aria-description={`Avatar for ${name}`}
      >
        <div>
          {i1}
          {i2 ?? ""}
        </div>
      </div>
    );
  }

  return (
    <div className={`${styles.wrapper} ${styles[`size-${size}`]}`}>
      <img src={src} alt={`Avatar for ${name}`} loading="lazy" />
    </div>
  );
};
