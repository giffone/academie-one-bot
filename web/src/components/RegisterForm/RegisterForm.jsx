import React, { useEffect, useState } from "react";
import "./RegisterForm.css";
import {mytelegram}  from "../../hooks/mytelegram.js";

const RegisterForm = () => {
  const [login, SetLogin] = useState();
  const [firstName, SetFirstName] = useState();
  const [lastName, SetLastName] = useState();
  const {tg} = mytelegram();

  useEffect(() => {
    tg.MainButton.setParams({
      text: "Send data",
    });
  }, []);

  useEffect(() => {
    if (!login || !firstName || lastName) {
      tg.MainButton.hide();
    } else {
      tg.MainButton.show();
    }
  }, [login, firstName, lastName]);

  const onChangeLogin = (e) => {
    SetLogin(e.target.value);
  };

  const onChangeFirstName = (e) => {
    SetFirstName(e.target.value);
  };

  const onChangeLastName = (e) => {
    SetLastName(e.target.value);
  };

  return (
    <div className={"regform"}>
      <h3>Input you personal data</h3>
      <input
        className={"input"}
        type="text"
        placeholder={"cadet login"}
        value={login}
        onChange={onChangeLogin}
      />
      <input
        className={"input"}
        type="text"
        placeholder={"first name"}
        value={firstName}
        onChange={onChangeFirstName}
      />
      <input
        className={"input"}
        type="text"
        placeholder={"last name"}
        value={lastName}
        onChange={onChangeLastName}
      />
    </div>
  );
};

export default RegisterForm;
