import { format } from "date-fns";

export const serialiseDate = date => format(date, "YYYY-DD-MM HH:mm:ss");

export const displayDate = date =>
  date.toLocaleString("en-UK", {
    weekday: "long",
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
