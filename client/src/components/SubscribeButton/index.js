import React, { useEffect, useState, useRef } from "react";
import axios from "axios";
import { Button, Toaster } from "@blueprintjs/core";
import uniqid from "uniqid";
import styles from "./index.module.scss";
import { serialiseJourney } from "../JourneyForm/data";

export const SubscribeButton = props => {
  const { journey, className, channel, setChannel } = props;
  const clickedState = useState(false);
  const setClicked = clickedState[1];
  const channelState = [channel, setChannel];
  let toaster = useRef(null);
  useHandleClick(journey, clickedState, channelState, toaster.current);
  return (
    <Button
      onClick={() => setClicked(true)}
      className={`${styles.submitButton} ${className}`}
    >
      <i className={`${styles.icon} far fa-bell`} />
      Notify Me!
      <Toaster ref={toaster} />
    </Button>
  );
};

const useHandleClick = (journey, clickedState, channelState, toaster) => {
  const [clicked, setClicked] = clickedState;
  useEffect(() => {
    if (!clicked) {
      return;
    }
    requestSubscription(journey, toaster, channelState, setClicked);
  });
};

async function requestSubscription(journey, toaster, channelState, setClicked) {
  const [channel, setChannel] = channelState;
  await axios.post("http://127.0.0.1:7891/subscribe", {
    ...serialiseJourney(journey),
    channel: channel
  });
  toaster.show(successToast);
  setClicked(false);
  setChannel(uniqid());
}

const successToast = {
  message:
    "You're successfully subscribed! We'll let you know which bus to catch.",
  intent: "success",
  icon: "endorsed",
  timeout: 2000
};
