import React, { useEffect, useState } from "react";
import axios from "axios";
import "./index.module.scss";
import { Suggest } from "@blueprintjs/select";
import { MenuItem } from "@blueprintjs/core";

export const RouteSelector = () => {
  const [routes, setRoutes] = useState([]);
  useEffect(() => FetchDataEffect(routes, setRoutes), [routes]);
  return (
    <Suggest
      inputValueRenderer={x => x}
      items={Object.keys(routes)}
      itemRenderer={renderRouteID}
      itemPredicate={doesQueryMatchRouteID}
      onItemSelect={item => console.log(item)}
    />
  );
};

const doesQueryMatchRouteID = (query, routeID) =>
  routeID.toLowerCase().indexOf(query.toLowerCase()) >= 0;

const renderRouteID = (routeID, { handleClick }) => (
  <MenuItem key={routeID} text={routeID} onClick={handleClick} />
);

const FetchDataEffect = (routes, setRoutes) => {
  if (routes.length !== 0) {
    return;
  }
  (async function fetchData() {
    const routes = await axios.get("http://127.0.0.1:7891/getStops");
    setRoutes(routes.data);
  })();
};
