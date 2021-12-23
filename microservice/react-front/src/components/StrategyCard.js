import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import CardActionArea from "@material-ui/core/CardActionArea";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import { Link } from 'react-router-dom';


const useStyles = makeStyles((muiBaseTheme) => ({
  card: {
    maxWidth: 700,
    maxHeight: 800,
    margin: "20px",
    transition: "0.3s",
    boxShadow: "0 8px 40px -12px rgba(0,0,0,0.3)",
    "&:hover": {
      boxShadow: "0 16px 70px -12.125px rgba(0,0,0,0.3)",
    },
    backgroundColor: "#d99e89",
    borderRadius: "10px",
  },
  card2: {
    maxWidth: 500,
    textAlign: "center",
    marginLeft: "150px",
  },
  media: {
    paddingTop: "56.25%",
  },
  content: {
    textAlign: "left",
    padding: muiBaseTheme.spacing(3),
  },
  divider: {
    margin: `${muiBaseTheme.spacing(3)}px 0`,
  },
  heading: {
    fontWeight: "bold",
  },
  subheading: {
    lineHeight: 1.8,
  },
  avatar: {
    display: "flexbox",
    border: "2px solid white",
    "&:not(:first-of-type)": {
      marginLeft: -muiBaseTheme.spacing(3),
    },
  },
}));

const StrategyCard = (props) => {
  const classes = useStyles();

  const newTo = {
    pathname: `/strategies/${props.ID}`,
    strategyName: props.StrategyName,
  }

  return (
    <Card className={classes.card} square={true}>
      {/* <CardActionArea component={Link} to={`/strategies/${props.ID}`}> */}
      <CardActionArea component={Link} to={newTo}>
        <CardContent className={classes.ctent}>
          <Typography
            color="textSecondary"
            gutterBottom
            style={{ paddingLeft: "1px", display: "flex", fontWeight: 1000 }}
          >
            strategyID: {props.ID}
          </Typography>
          <br />
          <br />
          <Typography variant="h2" component="h1" style={{ fontWeight: 500 }}>
            {props.StrategyName}
          </Typography>
          <br />
          <br />
          <Card className={classes.card2}>
            <Typography color="secondary" style={{ paddingLeft: "50%" }}>
              CreatedAt: {props.CreatedAt}
            </Typography>
            <Typography color="secondary" style={{ paddingLeft: "50%" }}>
              UpdatedAt: {props.UpdatedAt}
            </Typography>
            <Typography color="secondary" style={{ paddingLeft: "50%" }}>
              Schedule: {props.Schedule}
            </Typography>
          </Card>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default StrategyCard;
