import Pusher from "pusher-js";

export const PUSHER_EVENTS = {
  DEPARTURE_NOTIFICATION: "departureNotification"
};
const PUSHER_API_KEY = "66f6e62226c2a035a177";
const PUSHER_API_OPTIONS = {
  cluster: "eu"
};

export class PusherAPI {
  static pusherInstance = new Pusher(PUSHER_API_KEY, PUSHER_API_OPTIONS);

  static subscribe(channelName, eventName, callback) {
    const channel = this.pusherInstance.subscribe(channelName);
    channel.bind(eventName, callback);
  }
}
