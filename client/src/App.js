import React from "react";
import "./App.scss";
import { RouteSelector } from "./components/RouteSelector";

function App() {
  return (
    <div className="container">
      <h1 className="heading">Select a Journey</h1>
      <RouteSelector />
    </div>
  );
}

export default App;
