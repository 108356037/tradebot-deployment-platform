import { useParams } from "react-router";
import Demo from "./demo";
import { MaterialUIFormSubmit } from "./ScheduleForm";
import { useContext, useState, useEffect } from "react";
import AuthContext from "../store/auth-context";
import axios from "axios";
import ReactJson from "react-json-view";

const StrategyDetail = (props) => {
  const { strategyId } = useParams();
  const { userInfo, accessToken } = useContext(AuthContext);
  const [strategyInfo, setStrategyInfo] = useState({});
  const [ready, setReady] = useState(false);
  //const [strategyName, setStrategyName] = useState("");

  //const strategyName = props.location.strategyName

  useEffect(() => {
    async function GetStrategyDetail() {
      const userId = userInfo.sub;
      const response = await axios.get(
        `https://strategy.${process.env.REACT_APP_API_HOST}/v1/${userId}/strategies/id/${strategyId}`,
        {
          headers: {
            Authorization: "Bearer " + accessToken,
          },
        }
      );
      setStrategyInfo(response.data.info);
      //setStrategyName(response.data.info.strategy_name);
      setReady(true);
    }

    if (userInfo && accessToken) {
      GetStrategyDetail();
    }
  }, [userInfo, strategyId, accessToken]);

  return (
    <>
      {ready && (
        <div>
          <MaterialUIFormSubmit
            formName="Schedule your strategy"
            formDescription={`Target strategy: ${strategyInfo.strategy_name}`}
            uid={userInfo.sub}
            strategyName={strategyInfo.strategy_name}
          />
          <Demo />
          <ReactJson
            src={strategyInfo}
            theme="google"
            displayDataTypes={false}
            onEdit={false}
          />
        </div>
      )}
    </>
  );
};

export default StrategyDetail;
