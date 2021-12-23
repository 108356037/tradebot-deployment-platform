import { useEffect, useState, useContext } from "react";
import axios from "axios";
import AuthContext from "../store/auth-context";
import Button from "@material-ui/core/Button";
import CircularProgress from "@material-ui/core/CircularProgress";
import Loader from "react-loader-spinner";

const JupyterPage = (props) => {
  const [loading, setLoading] = useState(true);
  const [ready, setReady] = useState(false);
  const [createInProgress, setCreateInProgress] = useState(false);
  const [url, setUrl] = useState("");
  const { userInfo, accessToken } = useContext(AuthContext);

  const SetupJupyter = async () => {
    setCreateInProgress(true);
    axios.defaults.headers.post["Authorization"] = "Bearer " + accessToken;
    let res = await axios.post(
      `https://user-resource.${process.env.REACT_APP_API_HOST}/resources/${userInfo.sub}/jupyter`
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
            if (rs === `${userId}-jupyter`) {
              setReady(true);
              setUrl(`https://${userId}.jupyter.srv.algotrade.dev/playground`);
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
          title="jupyter-frame"
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
          <h1>Your jupyter workspace isn't setup yet</h1>
          <Button onClick={SetupJupyter} variant="contained" color="primary">
            Setup Jupyter
          </Button>
        </div>
      )}
      {createInProgress && <CircularProgress />}
    </>
  );
};

export default JupyterPage;
