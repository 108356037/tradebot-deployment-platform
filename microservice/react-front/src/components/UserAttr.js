//import { Auth } from "aws-amplify";
import { Button, Card, Typography } from "@material-ui/core";
import { useState } from "react";
import { Fragment } from "react";
import { useContext } from "react";
import AuthContext from "../store/auth-context";

const ShowUserAttr = (props) => {
  const authCtx = useContext(AuthContext);
  const [show, setShow] = useState(false);

  // https://cognito-idp.ap-southeast-1.amazonaws.com/ap-southeast-1_z2fADDubD/.well-known/jwks.json

  const showHandler = () => {
    console.log(authCtx);
    setShow(true);
  };

  const hideHandler = () => {
    setShow(false);
  };

  return (
    <Fragment>
      <Button variant="outlined" onClick={showHandler}>
        Show awsCredentials
      </Button>
      <Button variant="outlined" color="secondary" onClick={hideHandler}>
        Hide awsCredentials
      </Button>
      {show && (
        <Card>
          <Typography>{JSON.stringify(authCtx.userInfo)}</Typography>
          <Typography variant="h4">{authCtx.accessToken}</Typography>
          <Typography variant="h5">{authCtx.idToken}</Typography>
        </Card>
      )}
    </Fragment>
  );
};

export default ShowUserAttr;
