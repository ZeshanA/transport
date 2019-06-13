import _ from "lodash";
import { serialiseDate } from "../../lib/date";

export const applyChanges = (journey, setJourney, newProperties) => {
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

export const readableRouteID = routeID => routeID.split("_")[1];

export const getStopsAfter = (fromStop, allStops) => {
  if (!fromStop) {
    return [];
  }
  let fromStopIDIndex = allStops.findIndex(stop => stop.id === fromStop.id);
  return allStops.slice(fromStopIDIndex + 1);
};

export const doesQueryMatchItem = (query, item) => {
  return (
    String(item)
      .toLowerCase()
      .indexOf(String(query).toLowerCase()) >= 0
  );
};

export const doesQueryMatchStop = (query, stop) =>
  doesQueryMatchItem(query, stop.name);

export const serialiseJourney = journey => ({
  ...journey,
  directionID: parseInt(journey.directionID),
  fromStop: journey.fromStop ? journey.fromStop.id : null,
  toStop: journey.toStop ? journey.toStop.id : null,
  arrivalTime: serialiseDate(journey.arrivalTime)
});
