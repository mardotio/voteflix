import { type CSSProperties } from "react";

import {
  ConfusedFaceEmoji,
  FaceVomitingEmoji,
  FrowningFaceEmoji,
  GrimacingFaceEmoji,
  GrinningFaceWithBigEyesEmoji,
  NeutralFaceEmoji,
  PartyingFaceEmoji,
  PileOfPooEmoji,
  SlightlySmilingFaceEmoji,
  SmilingFaceWithHeartEyes,
  SmilingFaceWithSmilingEyesEmoji,
  ThinkingFaceEmoji,
} from "../Emoji";
import styles from "./RatingSlider.module.scss";

interface RatingSliderProps {
  value: number | null;
  onChange: (v: number) => void;
}

const EMOJIS = [
  PileOfPooEmoji,
  FaceVomitingEmoji,
  FrowningFaceEmoji,
  GrimacingFaceEmoji,
  ConfusedFaceEmoji,
  NeutralFaceEmoji,
  SlightlySmilingFaceEmoji,
  SmilingFaceWithSmilingEyesEmoji,
  GrinningFaceWithBigEyesEmoji,
  SmilingFaceWithHeartEyes,
  PartyingFaceEmoji,
];

export const RatingSlider = ({ onChange, value }: RatingSliderProps) => {
  const effectingRating = value === null ? -1 : value;
  const Emoji =
    effectingRating < 0 ? ThinkingFaceEmoji : EMOJIS[effectingRating];

  if (!Emoji) {
    throw new Error("Rating is out of range");
  }

  return (
    <div className={styles.main}>
      <Emoji size={36} />
      <div className={styles["reaction-slider"]}>
        <input
          style={
            {
              "--range-value": `${(effectingRating / 10) * 100}%`,
            } as CSSProperties
          }
          type="range"
          name="reaction"
          min={effectingRating >= 0 ? 0 : -1}
          max="10"
          list="tickmarks"
          value={effectingRating}
          onChange={(e) => {
            const val = Number(e.target.value);
            if (val >= 0 && val <= 10) {
              onChange(val);
            }
          }}
        />
        <datalist id="tickmarks">
          {[1, 5, 10].map((v) => (
            <option key={v} value={v} label={`${v}`} />
          ))}
        </datalist>
      </div>
    </div>
  );
};
