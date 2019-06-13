import { format, parse } from "date-fns";

export const serialiseDate = date => format(date, "YYYY-MM-DD HH:mm:ss");

export const displayDate = date =>
  date.toLocaleString("en-UK", {
    weekday: "long",
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });

export const displayTime = date =>
  parse(date).toLocaleString("en-UK", {
    hour: "2-digit",
    minute: "2-digit"
  });

export const currentTime = () =>
  new Date(
    new Date().toLocaleString("en-US", { timeZone: "America/New_York" })
  );
