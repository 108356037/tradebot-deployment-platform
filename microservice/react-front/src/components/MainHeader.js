import { NavLink } from "react-router-dom";
import classes from './MainHeader.module.css'

const MainHeader = () => {

  const checkActive = (match, location) => {
    //some additional logic to verify you are in the home URI
    if (!location) return false;
    const { pathname } = location;
    return pathname === "/";
  }

  return (
    <header className={classes.header}>
      <nav>
        <ul>
          <li>
            <NavLink activeClassName={classes.active} isActive={checkActive} to="/">home</NavLink>
          </li>
          <li>
            <NavLink activeClassName={classes.active} to="/playground">jupyter</NavLink>
          </li>
          <li>
            <NavLink activeClassName={classes.active} to="/c9">c9</NavLink>
          </li>
          <li>
            <NavLink activeClassName={classes.active} to="/strategies">strategies</NavLink>
          </li>
          <li>
            <NavLink activeClassName={classes.active} to="/botpage">deploy bot</NavLink>
          </li>
        </ul>
      </nav>
    </header>
  );
};

export default MainHeader;
