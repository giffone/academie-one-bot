import "./App.css";
import { useEffect } from "react";
import { mytelegram } from "./hooks/mytelegram";
import Header from "./components/Header/Header";
import { Route, Routes } from "react-router-dom";
import RegisterForm from "./components/RegisterForm/RegisterForm";
import Menu from "./components/Menu/Menu";

function App() {
  const { tg } = mytelegram();

  useEffect(() => {
    tg.ready();
  }, []);

  return (
    <div className="App">
      <Header />
      <Routes>
        <Route index element={<Menu />} />
        <Route path={"regform"} element={<RegisterForm />} />
      </Routes>
    </div>
  );
}

export default App;
