import type { ReactNode } from "react";
import { createPortal } from "react-dom";

import { ArrowIcon, CloseIcon } from "../Icon";
import styles from "./Drawer.module.scss";

export interface DrawerProps {
  isOpen: boolean;
  onClose: () => void;
  onBack?: () => void;
  children: ReactNode;
  className?: string;
  header?: ReactNode;
  height?: string;
}

export const Drawer = ({
  isOpen,
  children,
  onClose,
  onBack,
  className,
  header,
  height,
}: DrawerProps) => {
  return createPortal(
    <div
      className={`${styles.container} ${isOpen ? styles.open : styles.closed}`}
      aria-hidden={!isOpen}
    >
      <div className={styles.main} style={height ? { height } : undefined}>
        {isOpen && (
          <>
            <div className={styles.header}>
              <button onClick={onClose} className={styles["header-button"]}>
                <CloseIcon size={32} />
              </button>
              {(header || onBack) && (
                <div className={styles["header-left"]}>
                  {onBack && (
                    <button
                      onClick={onBack}
                      className={styles["header-button"]}
                    >
                      <ArrowIcon size={24} />
                    </button>
                  )}
                  {header && <h3>{header}</h3>}
                </div>
              )}
            </div>
            <div className={`${styles.content} ${className ?? ""}`}>
              {children}
            </div>
          </>
        )}
      </div>
    </div>,
    document.getElementById("drawer-root")!,
  );
};
