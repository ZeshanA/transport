import React, { useState } from "react";
import { Card, MenuItem, Elevation, FormGroup } from "@blueprintjs/core";
import { DateInput } from "@blueprintjs/datetime";
import { Suggest } from "@blueprintjs/select";
import _ from "lodash";

import { SubscribeButton } from "../SubscribeButton";

import styles from "./index.module.scss";
import { displayDate } from "../../lib/date";

const readableRouteID = routeID => routeID.split("_")[1];

const formatRoutes = routes =>
  _.reduce(
    routes,
    (acc, routeObj, routeID) => {
      return { ...acc, [readableRouteID(routeID)]: routeObj };
    },
    {}
  );

export const JourneyForm = ({ routes }) => {
  const [journey, setJourney] = useState({});
  const changeJourney = changes => applyChanges(journey, setJourney, changes);
  const formattedRoutes = formatRoutes(routes);
  let stopsForRoute = getStopsForRoute(journey, formattedRoutes);
  return (
    <Card elevation={Elevation.TWO} className={styles.formCard}>
      <form className={styles.form}>
        <FormGroup
          className={styles.formInput}
          label="Route"
          labelInfo="(required)"
        >
          <SuggestionInput
            inputProps={{ placeholder: "Type a Route ID..." }}
            items={Object.keys(formattedRoutes)}
            onItemSelect={routeID => changeJourney({ routeID })}
            selectedItem={journey.routeID}
          />
        </FormGroup>
        <FormGroup
          className={styles.formInput}
          label="Direction"
          labelInfo="(required)"
        >
          <SuggestionInput
            inputProps={{ placeholder: "Select a Direction ID" }}
            items={["0", "1"]}
            onItemSelect={directionID => changeJourney({ directionID })}
            selectedItem={journey.directionID}
          />
        </FormGroup>
        <FormGroup
          className={styles.formInput}
          label="From Stop"
          labelInfo="(required)"
        >
          <SuggestionInput
            inputProps={{ placeholder: "Select a source stop..." }}
            items={stopsForRoute}
            onItemSelect={fromStop => changeJourney({ fromStop })}
            itemPredicate={doesQueryMatchStop}
            selectedItem={journey.fromStop}
          />
        </FormGroup>
        <FormGroup
          className={styles.formInput}
          label="To Stop"
          labelInfo="(required)"
        >
          <SuggestionInput
            inputProps={{ placeholder: "Select a destination stop..." }}
            items={getStopsAfter(journey.fromStop, stopsForRoute)}
            itemRenderer={renderStopOption}
            onItemSelect={toStop => changeJourney({ toStop })}
            itemPredicate={doesQueryMatchStop}
            selectedItem={journey.toStop}
          />
        </FormGroup>
        <FormGroup
          className={`${styles.formInput} ${styles.wideInput}`}
          label="Desired Arrival Time"
          labelInfo="(required)"
        >
          <DateInput
            placeholder="Select when you'd like to arrive..."
            formatDate={displayDate}
            parseDate={Date.parse}
            defaultValue={new Date()}
            minDate={new Date()}
            timePrecision="minute"
          />
        </FormGroup>
        <SubscribeButton
          className={`${styles.formInput} ${styles.wideInput}`}
          journey={{ journey }}
        />
      </form>
    </Card>
  );
};

function getStopsForRoute(journey, routes) {
  let stopsForRoute = [];
  if (journey.routeID && journey.directionID) {
    stopsForRoute = routes[journey.routeID][journey.directionID];
  }
  return stopsForRoute;
}

const applyChanges = (journey, setJourney, newProperties) => {
  const resetOnChange = {
    routeID: ["fromStop", "toStop"],
    directionID: ["fromStop", "toStop"],
    fromStop: ["toStop"]
  };

  // Iterate over every changed property and construct an object
  // containing all properties that have actually changed, or need to be reset
  // as a result of a property further up the tree changing.
  let allChanges = _.reduce(
    newProperties,
    (allChanges, newValue, key) => {
      // If the property hasn't actually changed, move on
      if (journey[key] === newValue) {
        return allChanges;
      }
      // If there are no properties that need to be reset, we only have the one
      // property changing
      if (!resetOnChange[key]) {
        return { ...allChanges, [key]: newValue };
      }
      // Start with the an object containing the changed property and its new
      // value and then add all the other properties that now need to be reset
      // to null
      let changedProperties = resetOnChange[key].reduce(
        (acc, property) => ({ ...acc, [property]: null }),
        { [key]: newValue }
      );
      // Combine with the object containing changes made due to previous
      // properties
      return { ...allChanges, ...changedProperties };
    },
    {}
  );

  // Keep all properties of the existing journey object, with our changed
  // properties overwriting existing ones as needed
  setJourney({ ...journey, ...allChanges });
};

const SuggestionInput = props => {
  let propsToPass = {
    inputValueRenderer: item => (item.name ? item.name : item),
    itemRenderer: renderOption,
    noResults: <MenuItem disabled={true} text="No results." />,
    itemPredicate: doesQueryMatchItem,
    ...props
  };
  return <Suggest {...propsToPass} />;
};

const renderStopOption = (stop, itemProps) =>
  renderOption(stop, { ...itemProps, text: stop.name, key: stop.id });

const getStopsAfter = (fromStop, allStops) => {
  if (!fromStop) {
    return [];
  }
  let fromStopIDIndex = allStops.findIndex(stop => stop.id === fromStop.id);
  return allStops.slice(fromStopIDIndex + 1);
};

const doesQueryMatchItem = (query, item) => {
  return (
    String(item)
      .toLowerCase()
      .indexOf(String(query).toLowerCase()) >= 0
  );
};

const doesQueryMatchStop = (query, stop) =>
  doesQueryMatchItem(query, stop.name);

const renderOption = (item, { modifiers, handleClick, text, key }) => {
  return (
    <MenuItem
      active={modifiers.active}
      key={key ? key : item}
      text={text ? text : item}
      onClick={handleClick}
    />
  );
};
