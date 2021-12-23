import React, { useState, useEffect } from "react";
import { Auth } from "aws-amplify";
//import awsconfig from "../aws-exports";
import { useCookies } from "react-cookie";
import Amplify from "aws-amplify";
import { onAuthUIStateChange } from "@aws-amplify/ui-components";
import AuthContext from "./auth-context";
// import axios from "axios";

//Amplify.configure(awsconfig);

Amplify.configure({
  Auth: {
    identityPoolId: "ap-southeast-1:0673a1aa-f95a-4652-9b37-9aa0c635bfc2",
    region: "ap-southeast-1",
    userPoolId: "ap-southeast-1_z2fADDubD",
    userPoolWebClientId: "6mr5lbve47dbs0trvbuatqb3tm",
    mandatorySignIn: true,
    // cookieStorage: {
    //   domain: "127.0.0.1",
    //   path: "/",
    //   expires: 30,
    // },
    oauth: {},
  },
});

const AuthProvider = (props) => {
  const [user, setUser] = useState();
  const [accessToken, setAccToken] = useState("");
  const [idToken, setIdToken] = useState("");
  const [cookies, setCookie] = useCookies(["id_token"]);
  //const [fakeCurrentDate, setFakeCurrentDate] = useState(new Date());

  useEffect(() => {
    return onAuthUIStateChange(async (nextAuthState, authData) => {
      if (nextAuthState && authData) {
        setUser(await authData.attributes);
        setAccToken(await authData.signInUserSession.accessToken.jwtToken);
        setIdToken(await authData.signInUserSession.idToken.jwtToken);
        setCookie("id_token", authData.signInUserSession.idToken.jwtToken, {
          path: "/",
          secure: true,
          domain: "." + process.env.REACT_APP_APP_DOMAIN,
          maxAge: 1800,
        });

        // axios.defaults.headers.post["Authorization"] = "Bearer " + accessToken;
        // const userId = user.sub;
        // await axios.post(
        //   `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userId}/ns`
        // );
      }
    });
  }, [idToken, accessToken, user, cookies, setCookie]);

  // reload every 25 min to get a new cookie
  const MINUTE_MS = 60000;

  useEffect(() => {
    const interval = setInterval(async () => {
      const sess = await Auth.currentSession();
      var idToken = sess.getIdToken().getJwtToken();
      var accToken = sess.getAccessToken().getJwtToken();
      //console.log("Setting cookie");
      setCookie("id_token", idToken, {
        path: "/",
        secure: true,
        domain: "." + process.env.REACT_APP_APP_DOMAIN,
        maxAge: 1800,
      });
      setIdToken(idToken);
      setAccToken(accToken);
      console.log(
        "If you see this too frequently, it means the auth provider has some problem"
      );
    }, MINUTE_MS);
    return () => clearInterval(interval); // This represents the unmount function, in which you need to clear your interval to prevent memory leaks.
  }, [setCookie]);

  return (
    <AuthContext.Provider
      value={{
        userInfo: user,
        idToken: idToken,
        accessToken: accessToken,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
