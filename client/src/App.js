import React, { useEffect, useState } from "react";
import axios from "axios";
import "./App.scss";
import { JourneyForm } from "./components/JourneyForm";
import { PUSHER_EVENTS, PusherAPI } from "./api/pusher";
import uniqid from "uniqid";
import { DepartureAlert } from "./components/DepartureAlert";

function App() {
  const [routes, setRoutes] = useState([]);
  const [channel, setChannel] = useState(uniqid());
  const [notification, setNotification] = useState({});

  // Fetch routes data if needed
  useEffect(() => fetchDataEffect(routes, setRoutes), [routes]);
  useEffect(() => subscribeToPusherChannel(channel, setNotification), [
    channel
  ]);

  return (
    <div>
      <DepartureAlert
        notification={notification}
        setNotification={setNotification}
      />

      <div className="navbar">
        <h1>
          <i className="icon far fa-bus" />
          DelayGuardian
        </h1>
      </div>

      <section className="container">
        <h2 className="heading">
          <i className="icon far fa-map-marked" />
          Subscribe to Departure Notifications
        </h2>
        <JourneyForm
          routes={routes}
          channel={channel}
          setChannel={setChannel}
        />
      </section>
    </div>
  );
}

const fetchDataEffect = (routes, setRoutes) => {
  if (routes.length !== 0) {
    return;
  }
  (async function fetchData() {
    const routes = await axios.get("http://127.0.0.1:7891/getStops");
    setRoutes(routes.data);
  })();
};

const subscribeToPusherChannel = (channel, setNotification) => {
  PusherAPI.subscribe(
    channel,
    PUSHER_EVENTS.DEPARTURE_NOTIFICATION,
    notification => {
      setNotification({ ...notification, isOpen: true });
    }
  );
};

export default App;
