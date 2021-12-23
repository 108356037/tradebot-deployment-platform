import { useEffect, useState, useContext } from "react";
import AuthContext from "../store/auth-context";
import axios from "axios";
import Button from "@material-ui/core/Button";
import CircularProgress from "@material-ui/core/CircularProgress";
import Loader from "react-loader-spinner";

const C9Page = (props) => {
  const [loading, setLoading] = useState(true);
  const [ready, setReady] = useState(false);
  const [createInProgress, setCreateInProgress] = useState(false);
  const [url, setUrl] = useState("");
  const { userInfo, accessToken } = useContext(AuthContext);

  const SetupC9 = async () => {
    setCreateInProgress(true);
    axios.defaults.headers.post["Authorization"] = "Bearer " + accessToken;
    let res = await axios.post(
      `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userInfo.sub}/c9`
    );
    if (res.status === 201) {
      window.location.reload();
    } else {
      console.log(res);
    }
  };

  useEffect(() => {
    const fetchUrl = async () => {
      const userId = userInfo.sub;
      try {
        let res = await axios.get(
          `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userId}/helm`,
          {
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
              Authorization: "Bearer " + accessToken,
            },
          }
        );
        if (res.data.result !== "") {
          for (const rs of res.data.result) {
            if (rs === `${userId}-c9`) {
              setReady(true);
              setUrl(`https://${userId}.c9.srv.algotrade.dev/ide.html`);
              break;
            }
          }
        }
      } catch (error) {
        setReady(false);
      }
      setLoading(false);
    };
    if (userInfo && accessToken) {
      fetchUrl();
    }
  }, [userInfo, accessToken]);

  return (
    <>
      {loading && (
        <Loader
          type="Puff"
          color="#00BFFF"
          height={100}
          width={100}
          timeout={3000} //3 secs
        />
      )}
      {ready && (
        <iframe
          src={url}
          frameBorder="5px"
          marginWidth="0"
          title="c9-frame"
          scrolling="auto"
          style={{
            display: "flex",
            overflow: "hidden",
            height: "100%",
            width: "100%",
            position: "absolute",
          }}
          //allowFullScreen
        />
      )}
      {!ready && (
        <div>
          <h1>Your C9 workspace isn't setup yet</h1>
          <Button onClick={SetupC9} variant="contained" color="primary">
            Setup C9
          </Button>
        </div>
      )}
      {createInProgress && <CircularProgress />}
    </>
  );
};

export default C9Page;
