import type { ReactNode } from "react";
import { createPortal } from "react-dom";

import { CloseIcon } from "../Icon";
import styles from "./Drawer.module.scss";

export interface DrawerProps {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
  className?: string;
  header?: ReactNode;
}

export const Drawer = ({
  isOpen,
  children,
  onClose,
  className,
  header,
}: DrawerProps) => {
  return createPortal(
    <div
      className={`${styles.container} ${isOpen ? styles.open : styles.closed}`}
    >
      <div className={styles.main}>
        <div className={styles.header}>
          <button onClick={onClose} className={styles["close-button"]}>
            <CloseIcon size={32} />
          </button>
          {header && <h3>{header}</h3>}
        </div>
        <div className={`${styles.content} ${className ?? ""}`}>{children}</div>
      </div>
    </div>,
    document.getElementById("drawer-root")!,
  );
};
