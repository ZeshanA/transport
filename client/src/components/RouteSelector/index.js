import React, { useEffect, useState } from "react";
import axios from "axios";
import _ from "lodash";
import "./index.module.scss";

export const RouteSelector = () => {
  const [routes, setRoutes] = useState([]);
  useEffect(() => FetchDataEffect(routes, setRoutes), [routes]);
  return (
    <select>
      {_.map(routes, (routeObj, routeID) => (
        <option key={routeID}>{routeID}</option>
      ))}
    </select>
  );
};

const FetchDataEffect = (routes, setRoutes) => {
  if (routes.length !== 0) {
    return;
  }
  (async function fetchData() {
    const routes = await axios.get("http://127.0.0.1:7891/getStops");
    setRoutes(routes.data);
  })();
};
