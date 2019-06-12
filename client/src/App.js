import React, { useEffect, useState } from "react";
import axios from "axios";
import "./App.scss";
import { JourneyForm } from "./components/JourneyForm";

function App() {
  const [routes, setRoutes] = useState([]);

  // Fetch routes data if needed
  useEffect(() => FetchDataEffect(routes, setRoutes), [routes]);

  return (
    <div>
      <div className="navbar">
        <h1>
          <i className="icon far fa-bus" />
          DelayGuardian
        </h1>
      </div>
      <section className="container">
        <h2 className="heading">Select a Journey</h2>
        <JourneyForm routes={routes} />
      </section>
    </div>
  );
}

const FetchDataEffect = (routes, setRoutes) => {
  if (routes.length !== 0) {
    return;
  }
  (async function fetchData() {
    const routes = await axios.get("http://127.0.0.1:7891/getStops");
    setRoutes(routes.data);
  })();
};

export default App;
