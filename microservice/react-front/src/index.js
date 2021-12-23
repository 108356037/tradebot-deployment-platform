import { HashRouter } from "react-router-dom";
import AuthProvider from "./store/AuthProvider";
import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

ReactDOM.render(
  <HashRouter>
    <AuthProvider>
      <App />
    </AuthProvider>
  </HashRouter>,

  document.getElementById("root")
);
