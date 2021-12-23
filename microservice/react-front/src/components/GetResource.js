import Button from "@material-ui/core/Button";
import axios from "axios";
import { useContext, useState, useEffect } from "react";
import AuthContext from "../store/auth-context";

const ShowUserResource = (props) => {
  const { userInfo } = useContext(AuthContext);
  const [ready, setReady] = useState(false);

  async function GetResource() {
    const userId = userInfo.sub;
    const response = await axios.get(
      `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userId}/helm`
    );
    console.log(response);
  }

  useEffect(() => {
    if (userInfo) {
      setReady(true);
    }
  }, [userInfo]);

  return <>{ready && <Button onClick={GetResource}>Show resource</Button>}</>;
};

export default ShowUserResource;
