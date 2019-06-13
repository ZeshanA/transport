import React from "react";
import { Alert } from "@blueprintjs/core";
import { displayTime } from "../../lib/date";
import styles from "./index.module.scss";

export const DepartureAlert = ({ notification, setNotification }) => {
  if (!notification || !notification.vehicleID) {
    return null;
  }
  return (
    <Alert
      className={styles.alert}
      isOpen={notification.isOpen}
      icon="known-vehicle"
      onClose={() => setNotification({ ...notification, isOpen: false })}
    >
      <h2 className={styles.heading}>We found a perfect bus!</h2>
      <p className={styles.busMessage}>
        Take <strong>Bus {notification.vehicleID.split("_")[1]}</strong>{" "}
        arriving at{" "}
        <strong>{displayTime(notification.optimalDepartureTime)}</strong> to
        arrive at your destination by{" "}
        <strong>{displayTime(notification.predictedArrivalTime)}</strong>
      </p>
    </Alert>
  );
};
