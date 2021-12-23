import Button from "@material-ui/core/Button";
import { useCookies } from "react-cookie";
import { Auth } from "aws-amplify";

const SignOutBtn = (props) => {
  const [cookies, setCookie, removeCookie] = useCookies(["id_token"]);

  const SignOutHandler = async () => {
    removeCookie("id_token", {
      domain: ".srv.algotrade.dev",
    });
    await Auth.signOut();

    window.location.reload();
  };

  return <Button onClick={SignOutHandler}>SignOut</Button>;
};

export default SignOutBtn;
