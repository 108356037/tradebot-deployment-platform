import { useEffect, useState, useContext } from "react";
import Button from "@material-ui/core/Button";
import AuthContext from "../store/auth-context";
import axios from "axios";

const JupyterNewTabBtn = () => {
  const [ready, setReady] = useState(false);
  const [url, setUrl] = useState("");
  const { userInfo, accessToken } = useContext(AuthContext);

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
    };
    if (userInfo && accessToken) {
      fetchUrl();
    }
  }, [userInfo, accessToken]);

  const OpenJupyterInNewTab = () => {
    window.open(url, "JuypterLab");
  };

  return (
    <>
      {ready && (
        <Button
          variant="contained"
          color="secondary"
          onClick={OpenJupyterInNewTab}
        >
          Open Jupyter
        </Button>
      )}
    </>
  );
};

export default JupyterNewTabBtn;
