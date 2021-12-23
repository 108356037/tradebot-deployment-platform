import "./App.css";
import { withAuthenticator } from "@aws-amplify/ui-react";
import { Route, Switch, Redirect } from "react-router-dom";
import ShowUserAttr from "./components/UserAttr";
// import ShowUserResource from "./components/GetResource";
import CreateUserNamespace from "./components/CreateNamespace";
import JupyterPage from "./pages/JupyterPage";
import JupyterNewTabBtn from "./pages/OpenJupyterBtn";
import C9Page from "./pages/C9Page";
import SignOutBtn from "./components/Signout";
import MainHeader from "./components/MainHeader";
import Strategies from "./pages/Strategies";
import DeployBotPage from "./pages/DeployBotPage";
import StrategyDetail from "./components/StrategyDetail";

function App() {
  return (
    <div className="App">
      <MainHeader />
      <Switch>
        <Route path="/" exact>
          <ShowUserAttr />
          <SignOutBtn />
          <CreateUserNamespace />
        </Route>
        <Route path="/playground">
          <JupyterNewTabBtn />
          <JupyterPage />
        </Route>
        <Route path="/c9">
          <C9Page />
        </Route>
        <Route path="/strategies" exact>
          <Strategies />
        </Route>
        <Route path="/strategies/:strategyId" component={StrategyDetail} />
        <Route path="/botpage" exact>
          <DeployBotPage />
        </Route>
        {/* <StrategyDetail />
        </Route> */}
        <Redirect to="/" />
      </Switch>
    </div>
  );
}

export default withAuthenticator(App);
