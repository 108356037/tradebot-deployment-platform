import StrategyCard from "../components/StrategyCard";
import axios from "axios";
import { useState, useEffect, useContext } from "react";
import Loader from "react-loader-spinner";
import AuthContext from "../store/auth-context";


const Strategies = () => {
  const [isLoading, setLoading] = useState(true);
  const [userStrategies, SetUserStrategies] = useState([]);
  const { userInfo, accessToken } = useContext(AuthContext);

  useEffect(() => {
    async function ListStrategies() {
      const userId = userInfo.sub;
      const response = await axios.get(
        `https://strategy.${process.env.REACT_APP_API_HOST}/v1/${userId}/strategies`,
        {
          headers: {
            Authorization: "Bearer " + accessToken,
          },
        }
      );
      SetUserStrategies(response.data.info);
      setLoading(false);
    }
    if (userInfo && accessToken) {
      ListStrategies();
    }
  }, [userInfo, accessToken]);

  return (
    <>
      {isLoading && (
        <Loader type="Rings" color="skyblue" height={100} width={100} />
      )}
      {!isLoading && userStrategies.length === 0 && (
        <h1>No strategy deployed</h1>
      )}
      {!isLoading &&
        userStrategies.map((strategy, idx) => {
          return (
            <section key={idx}>
              <StrategyCard
                ID={strategy._id}
                StrategyName={strategy.strategy_name}
                CreatedAt={strategy.created_at}
                UpdatedAt={strategy.updated_at}
                Schedule={strategy.schedule}
              />
            </section>
          );
        })}
    </>
  );
};

export default Strategies;
