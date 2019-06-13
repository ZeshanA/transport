import React, { useState } from "react";
import { Card, MenuItem, Elevation, FormGroup } from "@blueprintjs/core";
import { DateInput } from "@blueprintjs/datetime";
import { Suggest } from "@blueprintjs/select";

import { SubscribeButton } from "../SubscribeButton";

import styles from "./index.module.scss";
import { currentTime, displayDate } from "../../lib/date";
import {
  applyChanges,
  doesQueryMatchItem,
  doesQueryMatchStop,
  getStopsAfter,
  readableRouteID
} from "./data";

export const JourneyForm = ({ routes, channel, setChannel }) => {
  const [journey, setJourney] = useState({ arrivalTime: currentTime() });
  const changeJourney = changes => applyChanges(journey, setJourney, changes);
  let stopsForRoute = getStopsForRoute(journey, routes);
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
            items={Object.keys(routes)}
            itemRenderer={(routeID, extra) =>
              renderOption(readableRouteID(routeID), extra)
            }
            inputValueRenderer={readableRouteID}
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
            itemRenderer={renderStopOption}
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
            placeholder="Select when you'd like to arrive.. ."
            formatDate={displayDate}
            parseDate={Date.parse}
            defaultValue={currentTime()}
            minDate={currentTime()}
            timePrecision="minute"
            onChange={arrivalTime => changeJourney({ arrivalTime })}
          />
        </FormGroup>

        <SubscribeButton
          className={`${styles.formInput} ${styles.wideInput}`}
          journey={journey}
          channel={channel}
          setChannel={setChannel}
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

const renderStopOption = (stop, itemProps) =>
  renderOption(stop, { ...itemProps, text: stop.name, key: stop.id });
