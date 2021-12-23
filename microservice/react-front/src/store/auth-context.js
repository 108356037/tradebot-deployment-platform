import React from "react";

const AuthContext = React.createContext({
  idToken: "",
  accessToken: "",
  userInfo: {},
  getIdToken: () => {},
  getAccessToken: () => {},
});

export default AuthContext;
