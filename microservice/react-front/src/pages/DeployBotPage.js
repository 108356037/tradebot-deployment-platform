import { Button, Icon, TextField, Paper, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import AuthContext from "../store/auth-context";
import JSONInput from "react-json-editor-ajrm";
import locale from "react-json-editor-ajrm/locale/en";
import axios from "axios";
import React, { useReducer, useContext, useEffect, useState } from "react";

const DeployBotPage = () => {
  const { userInfo } = useContext(AuthContext);
  const [postBody, SetPostBody] = useState({})
  useEffect(() => {
    if (userInfo) {
      console.log("userinfo get");
    }
  }, [userInfo]);

  const [formInput, setFormInput] = useReducer(
    (state, newState) => ({ ...state, ...newState }),
    {
      botname: "",
    }
  );

  const handleInput = (evt) => {
    const name = evt.target.name;
    const newValue = evt.target.value;
    setFormInput({ [name]: newValue });
    //console.log(formInput)
  };

  const handleSubmit = async (evt) => {
    evt.preventDefault();
    const uid = userInfo.sub
    const bottg = formInput.botname;
    console.log(postBody)

    const response = await axios.post(
      `https://tradebot.api.algotrade.dev/tradebot/${uid}/${bottg}`,
      postBody
    );
    console.log(response.data);
  };

  const useStyles = makeStyles((theme) => ({
    button: {
      margin: theme.spacing(1),
    },
    leftIcon: {
      marginRight: theme.spacing(1),
    },
    rightIcon: {
      marginLeft: theme.spacing(1),
    },
    iconSmall: {
      fontSize: 20,
    },
    root: {
      padding: theme.spacing(3, 2),
    },
    container: {
      display: "flex",
      flexWrap: "wrap",
    },
    textField: {
      marginLeft: theme.spacing(1),
      marginRight: theme.spacing(1),
      width: 400,
    },
  }));
  const classes = useStyles();

  return (
    <div>

      <JSONInput
        locale={locale}
        placeholder={{ "PUBLISHER": "the-publisher-name" }}
        confirmGood={false}
        onKeyPressUpdate={true}
        //theme={"light_mitsuketa_tribute"}
        style={{
          outerBox: { height: "auto", maxHeight: "350px", width: "100%" },
          container: {
            height: "auto",
            maxHeight: "350px",
            width: "100%",
            overflow: "scroll"
          },
          body: { minHeight: "45px", width: "100%" }
        }}
        // onChange={(v: Value) => onChange(v, props)}
        onChange={(v) => SetPostBody(v.json)} />
      <Paper className={classes.root}>
        <Typography variant="h5" component="h3">
          Deploy your Bot
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="Schedule"
            id="margin-normal"
            name="botname"
            defaultValue={formInput.botname}
            className={classes.textField}
            helperText="Fill in the botname here"
            onChange={handleInput}
          />
          {/* <TextField
            label="Name"
            id="margin-normal"
            name="name"
            defaultValue={formInput.email}
            className={classes.textField}
            helperText="Enter your full name"
            onChange={handleInput}
          />
          <TextField
            label="Email"
            id="margin-normal"
            name="email"
            defaultValue={formInput.name}
            className={classes.textField}
            helperText="e.g. name@gmail.com"
            onChange={handleInput}
          /> */}
          <Button
            type="submit"
            variant="contained"
            color="primary"
            className={classes.button}
          >
            Deploy<Icon className={classes.rightIcon}></Icon>
          </Button>
        </form>
      </Paper>
    </div>
  )
}

export default DeployBotPage