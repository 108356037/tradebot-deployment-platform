import React, { useReducer, useContext, useEffect, useState } from "react";
import { Button, Icon, TextField, Paper, Typography } from "@material-ui/core";
import { makeStyles } from "@material-ui/core/styles";
import axios from "axios";
import AuthContext from "../store/auth-context";
import SuccessSnackBar from "./SuccessSnackBar";

export function MaterialUIFormSubmit(props) {
  const { accessToken } = useContext(AuthContext);
  const [showSnackBar, setShowSnackBar] = useState(false);

  useEffect(() => {
    if (accessToken) {
      console.log("token initialized");
    }
  }, [accessToken]);

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

  // const [formInput, setFormInput] = useReducer(
  //   (state, newState) => ({ ...state, ...newState }),
  //   {
  //     name: "",
  //     email: ""
  //   }
  // );

  const [formInput, setFormInput] = useReducer(
    (state, newState) => ({ ...state, ...newState }),
    {
      schedule: "",
    }
  );

  const handleSubmit = async (evt) => {
    evt.preventDefault();
    const uid = props.uid;
    const strategyName = props.strategyName;
    let postData = {
      schedule: formInput.schedule,
    };

    let axiosConfig = {
      headers: {
        "Content-Type": "application/json;charset=UTF-8",
        Authorization: "Bearer " + accessToken,
      },
    };
    const response = await axios.patch(
      `https://strategy.api.algotrade.dev/v1/${uid}/strategies/${strategyName}/schedule`,
      postData,
      axiosConfig
    );
    console.log(response.data);
    if (response.data.status === "success") {
      setShowSnackBar(true);
    }
  };

  const handleInput = (evt) => {
    const name = evt.target.name;
    const newValue = evt.target.value;
    setFormInput({ [name]: newValue });
  };

  const classes = useStyles();

  return (
    <div>
      <Paper className={classes.root}>
        <Typography variant="h5" component="h3">
          {props.formName}
        </Typography>
        <Typography component="p">{props.formDescription}</Typography>

        <form onSubmit={handleSubmit}>
          <TextField
            label="Schedule"
            id="margin-normal"
            name="schedule"
            defaultValue={formInput.schedule}
            className={classes.textField}
            helperText="Fill in the crontab sechdule here"
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
            Schedule<Icon className={classes.rightIcon}></Icon>
          </Button>
        </form>
      </Paper>
      {showSnackBar && <SuccessSnackBar />}
    </div>
  );
}
