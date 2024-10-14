import React from "react";
import Button from "../Button/Button";
import { mytelegram } from "../../hooks/mytelegram";
import "./Header.css"

const Header = () => {
  const { user, onClose} = mytelegram()
  return (
    <div className={"header"}>
      <Button onClick={onClose}>Close</Button>
      <span className={"username"}>
        {user?.username}
      </span>
    </div>
  )
};

export default Header;
