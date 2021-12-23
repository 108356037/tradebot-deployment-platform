import Button from "@material-ui/core/Button";
import axios from "axios";
import { useEffect, useState, useContext } from "react";
import AuthContext from "../store/auth-context";
import CircularProgress from "@material-ui/core/CircularProgress";

const CreateUserNamespace = (props) => {
  const [ready, setReady] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const { userInfo, accessToken } = useContext(AuthContext);

  async function PostResource() {
    setIsLoading(true);
    const userId = userInfo.sub;
    axios.defaults.headers.post["Authorization"] = "Bearer " + accessToken;
    await axios.post(
      `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userId}/ns`
    );
    //console.log(response);
    setIsLoading(false);
    window.location.reload();
  }

  useEffect(() => {
    const checkRs = async () => {
      const userId = userInfo.sub;
      try {
        let res = await axios.get(
          `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userId}/helm`,
          {
            headers: {
              //"Content-Type": "application/x-www-form-urlencoded",
              Authorization: "Bearer " + accessToken,
            },
          }
        );
        if (res.data.result === "") {
          setReady(true);
        } else {
          setReady(false);
        }
      } catch (error) {
        console.log("error in fetching api, ", error);
      }
    };
    if (userInfo && accessToken) {
      checkRs();
    }
  }, [userInfo, accessToken]);

  return (
    <>
      {ready && (
        <Button onClick={PostResource} variant="contained" color="primary">
          Setup Workspace
        </Button>
      )}
      {isLoading && <CircularProgress />}
    </>
  );
};

export default CreateUserNamespace
