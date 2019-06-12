import React from "react";
import styles from "./index.module.scss";
import { Button } from "@blueprintjs/core";

export const SubscribeButton = ({ journey, className }) => {
  return (
    <Button
      onClick={handleClick}
      className={`${styles.submitButton} ${className}`}
    >
      <i className={`${styles.icon} far fa-bell`} />
      Notify Me!
    </Button>
  );
};

const handleClick = () => {
  alert("Clicked!");
};
